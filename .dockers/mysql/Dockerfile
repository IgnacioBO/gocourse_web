#Partiremos de la imagen de mysql con la version 5,7 https://hub.docker.com/_/mysql
FROM mysql:5.7

#Ejecutar las siguientes funcionaldiades para LEVANTAR LA BBDD
#Por defecto Docker usa /bin/sh en Linux
#Lo sigueinte cambia el sheel a bash, los siguientes comandos (como RUN) se ejecutara en bash en ves de sh
SHELL ["/bin/bash", "-c"]

#La levantaremos en el puerto 3306 (solo informativo, esto debe levantarse al levantar el container o en docker compose)
EXPOSE 3306
