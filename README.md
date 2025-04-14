# ResponsibleAPI-Go

## This is a personal project do not use it in production

- API
    - Development environment
    - docker build -f dev.Dockerfile . -t responsibleapi
    - ssh into the app container
        - docker exec -i -t responsible-api sh
    - ssh into the Mysql container
        - docker exec -i -t responsible-api-db sh