# Базовые команды
`go build .` - сборка бинаря
`go mod vendor` - подтянуть зависимости

# Структура проекта
./main.go - точка входа.

./pkg/bot.go - логика ответов бота.

./pkg/models/user.go - логика объекта Пользователь (например, внесение изменений в значения полей)

./pkg/store - логика работы с mongo db

# ENV
TELEGRAM_TOKEN=<bot token>

LOG_DIST=log.txt 

ADMINS="golanger brocmc"

# О работе на удаленке
`supervisord -c supervisor.conf -d` - запустить бота в режиме демона 
(только после сборки!)
