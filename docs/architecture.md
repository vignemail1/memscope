# Architecture cible

## Principes

Séparer acquisition, normalisation, analyse et présentation. Le cœur ne dépend pas de Windows, d'un pilote ou du format de sortie. Les collecteurs n'effectuent aucune recommandation. L'analyse ne réalise aucun accès matériel.

## Arborescence prévue

```text
cmd/memscope/
internal/cli/
internal/model/
internal/collect/inventory/
internal/collect/spd/
internal/collect/runtime/
internal/platform/windows/
internal/recommend/
internal/export/
internal/compare/
testdata/snapshots/
testdata/spd/
docs/
```

## Contrats de collecteurs

Chaque collecteur expose identité/version, capacités, plateformes prises en charge, privilèges nécessaires et collecte annulable via context.Context. Il retourne observations, diagnostics et erreurs structurées. Une erreur de source ne détruit pas les résultats des autres sources.

Prévoir InventoryProvider, SpdProvider et RuntimeMemoryProvider. Le décodage SPD reste indépendant du transport matériel. Les sources importées depuis un fichier portent source=import, jamais observed sans provenance d'origine conservée.

## Pipeline

1. Déterminer les capacités sans élévation automatique.
2. Collecter avec délais bornés et limites de taille.
3. Valider les valeurs et normaliser les unités.
4. Conserver les contradictions entre sources.
5. Produire un snapshot immuable.
6. Analyser ce snapshot avec un moteur de règles versionné.
7. Afficher ou exporter de manière déterministe.

Un snapshot est une acquisition sur un intervalle, pas nécessairement un instant matériel atomique. Chaque observation conserve son horodatage et sa fraîcheur.

## Accès privilégié

Si nécessaire, isoler l'accès matériel dans un composant minimal, signé et auditable. Pas de service installé implicitement. IPC avec ACL restrictives, validation des messages et opérations en liste blanche. Aucune API arbitraire de lecture/écriture mémoire physique. Les collectes sur un bus partagé sont sérialisées.

## Compatibilité

Utiliser les contraintes de compilation Go pour les collecteurs Windows. Les autres plateformes doivent retourner unsupported pour l'acquisition, mais pouvoir analyser un snapshot. Documenter séparément compatibilité du binaire Go, du composant natif, du CPU/chipset et du type mémoire.

## Extensions futures

Couverture CPU caches/instructions/topologie, SPD DDR4/DDR5, timings secondaires/tertiaires, télémétrie validée, benchmarks et imports de tests. Aucune feature n'est activée sur une plateforme non identifiée par simple ressemblance.
