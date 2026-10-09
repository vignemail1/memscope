# internal/model

Types internes partagés et validation indépendante de la plateforme. model.go reste compatible avec les contrats internal/collect. Ses commentaires initiaux sur la validation future sont remplacés fonctionnellement par les méthodes ajoutées dans value.go et validation.go.

## API

- NewText, NewUnsigned, NewBoolean : valeurs présentes, y compris chaîne vide, zéro et false.
- NewDecimal : rejette NaN et les infinis.
- Value.Validate, Format, Clone : union à une seule variante, format indépendant de la locale, copie indépendante.
- Status.Validate et Source.Validate : valeurs de provenance autorisées.
- Observation.Validate : identité, périmètre, paramètre, statut, source, timestamp et présence de valeur.
- Snapshot.Validate : métadonnées, intervalle, identités uniques, parents sans cycle, références et profils.
- ErrInvalid : compatible avec errors.Is sur erreurs contextualisées.

## Limites explicites

Ce package n'effectue aucun accès matériel. Il ne définit pas encore le codec JSON public ni une validation des versions de schéma. Les unités et bornes par paramètre seront vérifiées par une couche dédiée ; une unité vide peut être valide pour une donnée sans dimension. Les diagnostics et capacités ne font pas encore l'objet d'une validation métier.

Invalid et unavailable n'ont aucune valeur normalisée. Les observations restent datées dans l'intervalle de collecte ; les futurs imports historiques devront conserver leur propre intervalle, pas les présenter comme une nouvelle mesure.

Les types exportés restent mutables : après validation, l'appelant doit éviter les mutations concurrentes. Clone ne copie qu'une Value, pas tout le snapshot. Les recommandations constructeur restent distinctes des observations actuelles.

## Tests à exécuter

```sh
go test ./internal/model
go vet ./internal/model
```

Tests écrits mais non exécutés lors de la préparation de ce commit. Les contrôles de formatage et de compilation restent à réaliser dans la chaîne Go.
