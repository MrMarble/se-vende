# Se vende

Simple bot de telegram que extrae los datos de un anuncio de wallapop y los muestra en un mensaje de telegram. Pruebalo en https://t.me/sevende_bot

## Instalación

```sh
go build
```

## Uso

### Variables de entorno

- `TG_TOKEN`: Token de telegram. https://core.telegram.org/bots


### Docker

```sh
docker run -d --name sevende -e TG_TOKEN mrmarble/se-vende
```
