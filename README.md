# memscope

CLI de diagnostic matériel et mémoire pour Windows, développée en Go, visant progressivement une couverture comparable aux informations pertinentes de CPU-Z.

## État du projet

Phase de conception : ce dépôt contient les spécifications, pas encore un programme exécutable. Toutes les commandes ci-dessous sont prévues et non disponibles. La parité CPU-Z n'est pas acquise : chaque donnée dépendra des plateformes et collecteurs validés.

Version cible demandée : Go 1.27. Sa disponibilité et les versions des outils compatibles devront être vérifiées avant l'initialisation du module et de la CI. Ne pas utiliser implicitement une toolchain préliminaire.

## Fonctionnalités prévues

- Inventaire CPU, carte mère, BIOS et modules mémoire.
- Lecture des données SPD et profils JEDEC/XMP/EXPO, selon support matériel.
- Lecture du débit, des timings et des tensions actifs lorsque disponibles.
- Comparaison entre profils disponibles et configuration active.
- Recommandations BIOS documentées, distinguant profils constructeur et candidats expérimentaux.
- Export CSV des valeurs actuelles et des profils ; snapshots JSON versionnés.
- Comparaison avant/après et rattachement de résultats de tests.
- Diagnostic des capacités, permissions et limitations des collecteurs.

## Usage prévu

```powershell
memscope inspect
memscope memory current
memscope memory profiles
memscope recommend
memscope doctor
memscope export --scope current --format csv --output current.csv
memscope export --scope profiles --format csv --output profiles.csv
memscope snapshot --output before.json
memscope compare before.json after.json
```

Aucune commande ne modifiera le BIOS, les tensions, les timings ou le SPD. Les changements proposés seront appliqués manuellement dans l'UEFI puis validés par des tests. Les valeurs inconnues ne seront jamais remplacées par zéro ou par les valeurs d'un profil.

## Plateformes

Cible initiale : Windows 11 x64. Les analyses de snapshots et les exports devront rester portables. Le développement et la compilation du composant Go pourront être réalisés sous macOS ; les composants natifs et pilotes nécessiteront une chaîne Windows distincte. Windows ARM64, Linux et macOS comme sources matérielles ne sont pas promis dans la première version.

## Documentation

- [Index et décisions](docs/README.md)
- [Architecture](docs/architecture.md)
- [Accès aux informations](docs/hardware-access.md)
- [Contrat CLI](docs/cli.md)
- [Données et exports](docs/data-model.md)
- [Politique de recommandations](docs/recommendations.md)
- [Features et roadmap](docs/roadmap.md)
- [Contribution et validation](CONTRIBUTING.md)
- [Sécurité](SECURITY.md)

## Limites

Un profil présent n'est pas forcément activé ou compatible avec toute la plateforme. Un délai CAS calculé n'est pas la latence CPU-vers-RAM. Une tension déclarée par un profil n'est pas une mesure. Un démarrage réussi ne prouve pas la stabilité. Les meilleurs réglages ne peuvent pas être déduits universellement du SPD.

## Licence

Licence du projet à décider par le propriétaire avant redistribution. Aucune licence de SDK ou pilote tiers n'est présumée autoriser son intégration.
