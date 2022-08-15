# FPSBot

Voici le repo git pour le bot Discord du serveur FPSB.

## Installation

Tout d'abord, clonez le repo sur votre machine.

Il faudra ensuite créer un fichier `.env` à la racine du dossier, contenant les variables suivantes :

- DISCORD_TOKEN : le token du bot
- DISCORD_GUILD : l'ID du serveur
- VOICE_CATEGORY : l'ID de la catégorie des salons vocaux
- VOICE_CHANNEL : l'ID du salon servant à la création des salons vocaux temporaires

Il suffira ensuite de compiler le programme avec la commande `go build main.go` et de lancer le bot avec la commande `./main`.