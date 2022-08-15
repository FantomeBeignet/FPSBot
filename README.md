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

## Fonctionnalités

### Commandes de base

Le bot dispose de quelques commandes de base :

- /v2, /v3, /v4 et /v5 renvoient le lien de la verision du forum correspondante.
- /ordre renvoie une image de l'ordre d'apprentissage recommandé pour les débutants.

### Autres fonctionnalités

Ce bot permettra la création de salons vocaux temporaires quand des utilisateurs souhaitent discuter hors des salons principaux.

Rejoindre le salon dont l'ID sera fourni dans le fichier `.env` créera un nouveau salon au nom de l'utilisateur et le déplacera automatiquement dedans. Ce salon sera supprimé automatiquement une fois vide.
