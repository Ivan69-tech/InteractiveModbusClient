## Contexte du projet

Il est souvent très utile de surveiller un serveur modbus dans sa globalité dans le cadre de la gestion des interfaces dans un projet industriel.
Ce projet vise à donner une facon simple de requêter un serveur modbus sur de nombreux registres en mêmes temps.

- [Lien du projet Github](https://github.com/Ivan69-tech/gotools)

## Inputs du programme

------ A. Le programme prend en entrée un tableau CSV de la forme suivante (exemple) :
La première ligne ne doit pas être indiqué. Les espaces ne sont pas supportés. Voir les exemples sur le repo.


| Nom du signal                       | Registre           | Taille de la donnée | Position du bit | Type de donnée |
| ----------------------------------- | ------------------ | ------------------- | ----------------|----------------|
| Target state                        |        0xd102      |        uint32       |      100        |      holding   |
| Puissance mesurée                   |        9488        |        int32        |      100        |      input     |
| Fire Alarm                          |        0x4000      |        uint32       |      3          |      input     |


1. Registre --> l'hexadécimal et le décimal sont supportés
2. taille de la donnée --> uint16 / int16 / uint32 / int32 sont supportés
3. Position du bit --> si ce chiffre est supérieur à 32, le registre entier est retourné. Sinon 0 ou 1 est retourné.
4. Type de donnée --> holding / input / coil sont supportés

Si un signal binaire est lu à 1, la case correspondante devient rouge pour faciliter la lecture des signaux actifs.

/!\ Le tableau CSV doit être placé à la racine du binaire ou du .exe. IL DOIT SE NOMMER conf.csv.

------ B. il faut renseigner le host du serveur modbus sur la page d'accueil.

## Logs

Si vous rencontrez des problèmes lors de la lecture modbus, analyser les logs via le bouton "voir les logs"
Si une requête ModBus n'a pas pu aboutir, vous lirez 404 
Les logs sont limités à 70 lignes.

## Futures améliorations

1. Supporter le slave ID.
2. Améliorer la gestion des erreurs.
3. Améliorer la manière de renseigner le tableau des requêtes. Ce n'est pas très intuitif.
4. Héberger le projet sur un serveur distant type raspberry Pi (Attention au forward de l'host).
5. Possibilité de tracer les graphes des chaque signal en fonction du temps
6. Reprendre toute l'architecture du code en utilisant vue.js (faire le point 6 puis le point 5)
