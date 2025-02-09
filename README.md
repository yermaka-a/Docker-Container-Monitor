# Описание проекта: Docker Container Monitor 🌐🐳

## Оглавление

1. [Общая информация 📑](#общая-информация-)
2. [Запуск проекта  в **Docker** 🐋](#запуск-проекта-в-docker-)
3. [Архитектура проекта 🏢](#архитектура-проекта-)
   - [Сервисы](#сервисы)
4. [Взаимодействие сервисов 🌐](#взаимодействие-сервисов-)
5. [Настройки по умолчанию 🟢](#настройки-по-умолчанию-env-)
6. [Заключение ✅](#заключение-)

## Общая информация 📑
Проект "Docker Container Monitor" представляет собой систему мониторинга Docker-контейнеров, реализованную с использованием современных технологий, таких как Go, PostgreSQL, Nginx, RabbitMQ и React с TypeScript. Все четыре основных сервиса работают внутри одной Docker-сети, управляемой с помощью `docker-compose`, что обеспечивает простоту развертывания и масштабирования.

## Архитектура проекта 🏢
### Сервисы
1. **Backend (Go) - RESTful API**

- Этот сервис реализует RESTful API с двумя основными эндпоинтами: **GET** и **POST**.
- **GET/containers**: Возвращает массив контейнеров, содержащий:
  - `id` контейнера в базе данных
  - `ip` адрес контейнера внутри Docker-сети
  - `timeMs` - время последнего пинга
  - `pingDate` - дата последнего успешного пинга
- **POST/containers**: Позволяет отправлять массив контейнеров вручную (например, через Postman) для добавления или обновления данных в базе данных PostgreSQL.
- Доступ к API осуществляется через обратный прокси-сервер Nginx, который работает на порту, по умолчанию: **3001** (**http://localhost:3001**).

2. **Pinger (Go) - Асинхронный пингер**

- Этот сервис периодически (по умолчанию каждые 7 секунд) пингует Docker-контейнеры внутри сети, используя Docker API для Go.
- Для доступа к Docker-контейнерам используется volume: - ```/var/run/docker.sock:/var/run/docker.sock.```
- После пинга контейнеров, сервис отправляет результаты через сервис очередй **RabbitMQ** в backend, который сохраняет или обновляет данные в PostgreSQL.

3. **Frontend (React + TypeScript)**
- Интерфейс пользователя, реализованный с помощью React и TypeScript, отображает таблицу с данными контейнеров (id, ip, timeMs, pingDate).
- Он делает GET-запросы к backend через прокси Nginx, получая обновленные данные каждые несколько секунд.
- Frontend также работает на отдельном Nginx, который предоставляет статические файлы. 
- Доступ к frontend - (**http://localhost:5137**)(**Порт по умолчанию: 5137**)

4. **PostgreSQL**

- Сервис базы данных, который хранит информацию о контейнерах, полученную от backend.
- Обеспечивает надежное и эффективное хранение данных.

### Структура проекта
<details>
 <summary>🔍 Нажмите чтобы раскрыть</summary>

<ol>
    <li><strong>root</strong>
        <ul>
            <li>📁 Корневая папка: <code>root/</code>
                <ul>
                    <li>📁 back</li>
                    <li>📁 front</li>
                    <li>📁 nginx</li>
                    <li>📁 pinger</li>
                    <li>🗑️ .dockerignore</li>
                    <li>⚙️ .env</li>
                    <li>🗑️ .gitignore</li>
                    <li>📦 docker-compose.yaml</li>
                    <li>📜 README.md</li>
                </ul>
            </li>
        </ul>
    </li>
    <li><strong>backend</strong>
        <ul>
            <li>📁 back
                <ul>
                    <li>📁 internal
                        <ul>
                            <li>📁 app
                                <ul>
                                    <li>📄 app.go</li>
                                </ul>
                            </li>
                            <li>📁 config
                                <ul>
                                    <li>📄 config.go</li>
                                </ul>
                            </li>
                            <li>📁 db
                                <ul>
                                    <li>📄 db.go</li>
                                    <li>📄 models.go</li>
                                </ul>
                            </li>
                            <li>📁 services
                                <ul>
                                    <li>📄 connContainersRMQ.go</li>
                                    <li>📄 getContainers.go</li>
                                    <li>📄 updateContainers.go</li>
                                </ul>
                            </li>
                        </ul>
                    </li>
                    <li>🐳 Dockerfile</li>
                    <li>📄 main.go</li>
                </ul>
            </li>
        </ul>
    </li>
    <li><strong>frontend</strong>
        <ul>
            <li>📁 front
                <ul>
                    <li>📁 src
                        <ul>
                            <li>📁 components
                                <ul>
                                    <li>📁 Header
                                        <ul>
                                            <li>📄 Header.tsx</li>
                                        </ul>
                                    </li>
                                    <li>📁 Table
                                        <ul>
                                            <li>📄 Table.tsx</li>
                                            <li>📄 types.ts</li>
                                        </ul>
                                    </li>
                                </ul>
                            </li>
                            <li>📄 App.tsx</li>
                            <li>📄 index.css</li>
                            <li>📄 main.tsx</li>
                        </ul>
                    </li>
                    <li>🐳 Dockerfile</li>
                    <li>📄 index.html</li>
                    <li>📄 nginx.conf</li>
                </ul>
            </li>
        </ul>
    </li>
    <li><strong>proxy</strong>
        <ul>
            <li>📁 nginx
                <ul>
                    <li>📄 default.conf</li>
                    <li>🐳 Dockerfile</li>
                </ul>
            </li>
        </ul>
    </li>
    <li><strong>pinger</strong>
        <ul>
            <li>📁 pinger
                <ul>
                    <li>📁 internal
                        <ul>
                            <li>📁 app
                                <ul>
                                    <li>📄 app.go</li>
                                </ul>
                            </li>
                            <li>📁 config
                                <ul>
                                    <li>📄 config.go</li>
                                </ul>
                            </li>
                            <li>📁 models
                                <ul>
                                    <li>📄 models.go</li>
                                </ul>
                            </li>
                            <li>📁 services
                                <ul>
                                    <li>📄 services.go</li>
                                </ul>
                            </li>
                        </ul>
                    </li>
                    <li>🐳 Dockerfile</li>
                    <li>📄 main.go</li>
                </ul>
            </li>
        </ul>
    </li>
</ol>

</details>

### Схема работы ⟲

<details>
 <summary>🔍 Нажмите чтобы раскрыть</summary>
<img alt="Примерная схема работы проекта" src="https://github.com/user-attachments/assets/7e167f73-f167-4c98-80ec-d16198de16db" />
</details>

## Взаимодействие сервисов 🌐
- **Nginx** выступает в роли обратного прокси-сервера, обеспечивая доступ к backend через внешний порт 3001.
- **Pinger** асинхронно пингует контейнеры и отправляет результаты в backend через RabbitMQ.
- **Backend** обрабатывает данные, полученные от pinger, и сохраняет их в - *PostgreSQL*.
- **Frontend** запрашивает данные у backend и отображает их пользователю.

## Настройки по умолчанию: **.env** 🟢
Все настройки по умолчанию для портов, времени опроса контейнеров и т.д. Можно посмотреть в файле **.env** в корне проекта.

## Запуск проекта в **Docker** 🐋
Для запуска проекта используйте docker-compose для создания и управления контейнерами. Убедитесь, что все зависимости установлены и конфигурации настроены правильно.
1. Убедитесь, что Docker установлен
2. Скачайте архив с проектом или клонируйте репозиторий командой `git clone` по ссылке: [**Docker-Container-Monitor**](https://github.com/yermaka-a/Docker-Container-Monitor.git)
    
    ```bash
        git clone https://github.com/yermaka-a/Docker-Container-Monitor.git
    ```
3. Перейдите в директорию с проектом, если необходимо:
    ```bash
        cd название_директории
    ```
4. В терминале выполните следующую команду для запуска всех сервисов:
    ```bash
        docker-compose up --build
    ```
5. Проверьте работу сервисов:
    - После успешного запуска, вы сможете получить доступ к интерфейсу пользователя по адресу: ```http://localhost:5137/```.

6. Тестирование API:
    - Используйте Postman или любой другой инструмент для тестирования API, чтобы убедиться что эндопинты работают корректно.

## Заключение ✅
Проект "Docker Container Monitor" предоставляет мощный инструмент для мониторинга состояния Docker-контейнеров в реальном времени. Используя современные технологии, такие как **Go**, **PostgreSQL**, **Nginx**, **RabbitMQ** и **React** с **TypeScript**, было создано надежное и масштабируемое приложение.
