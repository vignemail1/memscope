# Contribution

## État initial

Documentation uniquement : aucun module Go, build, CI ou exécutable n'est encore présent. Ne pas annoncer que les commandes suivantes ont déjà été exécutées.

## Chaîne de développement cible

Go 1.27 est demandé ; confirmer sa disponibilité avant création du go.mod. Prévoir mise.toml avec versions figées, tâches fmt-check, vet, lint, test, fuzz, vuln et build. La sélection des outils de lint compatibles est une tâche préalable.

Contrôles prévus après ajout du code :

```sh
gofmt -l .
go vet ./...
go test ./...
```

La CI échoue si gofmt liste un fichier. Ajouter analyse des vulnérabilités, tests du CLI et builds Windows. CGO_ENABLED=0 est un objectif pour le cœur pur Go, pas une obligation mensongère pour un backend SDK natif. Les DLL/pilotes ont leur chaîne et validation séparées.

## Exigences de code

Erreurs contextualisées et typées, contextes annulables, délais bornés, absence de panic sur données externes, aucun nombre magique matériel non sourcé. Pas de logique d'accès matériel dans le moteur de recommandations. Pas de logs sur stdout structuré.

Tests : unités, golden CSV/JSON, fuzz SMBIOS/SPD/snapshots, limites de tailles, timeouts, erreurs partielles, conflits de sources et unités. Les tests matériels sont opt-in et ne doivent pas nécessiter de désactiver une protection Windows.

## Revue

Chaque PR indique objectif, feature Fxxx, limites, tests effectués, plateformes et impacts sécurité/licence. Ne pas prétendre valider un matériel non testé. Ajouter des fixtures anonymisées avec provenance et droits de redistribution.

## Documentation

Français pour la documentation initiale, identifiants et code en anglais. Maintenir synchronisés aide CLI, schémas et exemples. Marquer explicitement planned/experimental/supported. Toute dépendance native ou règle de tension nécessite un ADR et une revue spécifique.
