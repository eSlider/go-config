# Конфиг в Go: библиотек много, «единого решения» нет

*Что регулярно ломается в реальных сервисах, когда надо совместить YAML, `.env`, переменные окружения и вложенный `Config`.*

---

Вот абсолютно бытовая ситуация.

Есть `config.yaml` для локалки. Есть `.env.example`, который у каждого чуть свой. В проде значения прилетают через Docker/Kubernetes/systemd. В коде живет нормальный вложенный `Config`, а не плоская простыня.

И вот в этот момент становится ясно: в Go нет одного «очевидного» инструмента, который без плясок закрывает всю цепочку целиком.

Это не наезд на экосистему. В ней много сильных библиотек. Проблема в другом: почти каждая решает свой кусок, а швы между кусками остаются на команде.

---

## Проблема в одном примере

В коде:

```go
type Config struct {
    HTTP struct {
        Listen string
        TLS    bool
    }
    Database struct {
        URL string
    }
}
```

В `docker-compose.yml` — `http.listen`.
В Kubernetes — `HTTP_LISTEN`.
В YAML кто-то пишет `database.url`.
Провайдер PostgreSQL в документации советует `DATABASE_URL`.

Все правы. Но бинарнику от этого не легче.

Ему нужно:

1. Прочитать разные форматы.
2. Слить источники в понятном порядке приоритетов.
3. Заполнить типизированную структуру без ручного ада из `os.Getenv` и `strconv`.

Именно на этом месте «конфиг» перестает быть одной задачей и разваливается на три.

---

## Что я называю «единым» решением

Для себя я держу простой чек-лист:

- **Парсинг форматов** — YAML/JSON/INI/dotenv.
- **Слои и приоритеты** — defaults < repo config < local override < process env (по логике [12-factor](https://12factor.net/config)).
- **Нормализация ключей** — чтобы `sub-service`, `sub_service` и `SUB_SERVICE` жили в одном мире.
- **Декод в структуры** — конечная цель это `Config`, а не `map`.
- **Обратная кодировка** — уметь вывести эффективный конфиг обратно в файл.
- **CLI для операционки** — чтобы не писать вспомогательные `cmd/*` для каждой мелочи.
- **Тестируемость** — без скрытых глобалов, с фиксированным и проверяемым merge-поведением.

Важно: это не значит «одна библиотека обязана уметь всё». Нормально, когда инструмент закрывает 2-3 пункта. Ненормально, когда в README обещается «полный цикл», а сложные кейсы остаются «догадайтесь сами».

---

## Кратко по инструментам

### Viper

[Viper](https://github.com/spf13/viper) — первый выбор у многих. Большое сообщество, много примеров, привычный API.

Где боль: вложенные env-ключи и связка `AutomaticEnv` + `SetEnvKeyReplacer` + `BindEnv`. Проблема известная и давняя, это видно по тредам вроде [#641](https://github.com/spf13/viper/issues/641) и [#2001](https://github.com/spf13/viper/issues/2001).

Итог: рабочий вариант, особенно если команда уже на нем. Но с вложенным конфигом и сложным env-layout нужна дисциплина.

### Koanf

[Koanf](https://github.com/knadh/koanf) обычно воспринимается как более аккуратная композиция: providers, parsers, явный merge-порядок.

Плюс: пайплайн прозрачен.
Минус: часть решений все равно на вас (нормализация ключей, соглашения по env, стратегия декода).

### Env-first библиотеки

[caarlos0/env](https://github.com/caarlos0/env), [envconfig](https://github.com/kelseyhightower/envconfig), [cleanenv](https://github.com/ilyakaznacheev/cleanenv) отлично подходят, когда источник истины — env, а задача — быстро собрать типизированный `Config`.

Если же у вас YAML + env + несколько слоев, они не дадут весь конвейер «из коробки». Нужен клей.

### Dotenv-парсеры

[joho/godotenv](https://github.com/joho/godotenv) делает ровно то, что заявлено: корректно читает `.env`.

Это хороший кирпич. Но не целый дом.

### «Просто парсеры»

[encoding/json](https://pkg.go.dev/encoding/json), [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml), INI-библиотеки — хорошие парсеры.

Но они не решают сами по себе:
- порядок слоев,
- env-override,
- нормализацию имен ключей между форматами.

### mapstructure

[go-viper/mapstructure](https://github.com/go-viper/mapstructure) (v2) — по факту стандартный мост из `map[string]any` в структуру.

Это не парсер и не merge-движок. Его задача — декод.
Поэтому без аккуратной «середины» (о ней ниже) магии не будет.

---

## Почему вложенные структуры — главный тест

На плоском конфиге почти всё выглядит красиво. Проблемы приходят, когда структура становится реальной:

- **Embedding / `squash`**: где-то поля должны «подниматься», где-то жить в поддереве.
- **Несколько тегов на одно поле**: `json`, `yaml`, `mapstructure`, иногда `env`.
- **Списки в env**: `a,b,c`, JSON-строка, индексные ключи — у всех свои правила.
- **Слабая типизация**: особенно заметно на стыке `any`, `float64`, yaml-особенностей и env-строк.
- **`*Section` vs value**: отсутствие ключа, пустое значение и `nil` — не одно и то же.
- **«JSON в env»**: рабочий костыль, но часто больной в эксплуатации (кавычки, экранирование, логирование).

Если библиотека шикарна на плоском env, но сыпется на вложенных деревьях — это не «плохая библиотека». Просто ее зона оптимизации другая.

---

## Недостающая середина: map в центре пайплайна

Практически везде рабочая схема выглядит так:

```text
bytes -> nested map[string]any -> mapstructure -> struct
```

Левая часть — чтение источников.
Правая — декод в структуру.
А вот середину (merge + нормализация ключей) команды часто собирают сами и по-разному.

Почему это важно явно оформить:

- **Одна точка для кросс-форматной эквивалентности** (`sub-service` == `SUB_SERVICE` после нормализации).
- **Предсказуемые тесты** (можно проверять merged-map до декода).
- **Простой CLI** (`convert`, `merge`, `get` — это по сути операции вокруг той же map).

Это не призыв «всё переписать на map». Это призыв честно назвать центральный этап, от которого зависит поведение всей системы.

---

## Практический вывод

Универсальной «серебряной пули» нет.
Есть осознанный выбор того, **каким слоем вы управляете сами**, а что делегируете библиотеке.

Два правила, которые реально экономят время:

1. **Сразу зафиксируйте модель приоритетов и именования ключей.** Не «когда начнет гореть», а в первый день.
2. **Сделайте merge + normalizer отдельным, тестируемым слоем.** Это чаще окупается сильнее, чем замена одной библиотеки на другую.

---

## Где здесь `go-config`

[go-config](https://github.com/eslider/go-config) — попытка сделать именно эту «середину» предсказуемой:

- одинаковый `Codec`-подход для `env`, `yaml`, `json`, `ini`;
- deep merge для map и last-write-wins для скаляров;
- нормализация ключей через `LowerAlnum`;
- CLI `envc` для `convert / get / merge`.

Для контейнеров отдельный плюс: пакет `env` умеет слоить dotenv-файлы и затем применять **`WithCurrentEnvironment()`**. Это берет текущее окружение процесса (`os.Environ()` на момент `Map`) как верхний слой, поэтому одна и та же схема работает и локально, и в Docker/Kubernetes.

Это не «единственно правильный путь». Это одна из рабочих реализаций подхода, описанного выше.

Подробности API — в [README](https://github.com/eSlider/go-config/blob/main/README.md).
Архитектурные решения — в [ASR](https://github.com/eSlider/go-config/blob/main/docs/asr/README.md).

---

## Ссылки

- 12-factor — Config: <https://12factor.net/config>
- Awesome Go — Configuration: <https://awesome-go.com/configuration/>
- Viper: <https://github.com/spf13/viper>
- Viper issue #641 — nested env binding: <https://github.com/spf13/viper/issues/641>
- Viper issue #2001 — nested struct decode from env: <https://github.com/spf13/viper/issues/2001>
- Koanf: <https://github.com/knadh/koanf>
- caarlos0/env: <https://github.com/caarlos0/env>
- kelseyhightower/envconfig: <https://github.com/kelseyhightower/envconfig>
- cleanenv: <https://github.com/ilyakaznacheev/cleanenv>
- joho/godotenv: <https://github.com/joho/godotenv>
- go-viper/mapstructure (v2): <https://github.com/go-viper/mapstructure>
- encoding/json (stdlib): <https://pkg.go.dev/encoding/json>
- gopkg.in/yaml.v3: <https://github.com/go-yaml/yaml>
- Figment (Rust): <https://docs.rs/figment/latest/figment/>
- go-config — README: <https://github.com/eSlider/go-config/blob/main/README.md>
- go-config — ASR index: <https://github.com/eSlider/go-config/blob/main/docs/asr/README.md>
