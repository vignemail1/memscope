# Documentation de développement

Ces documents constituent un contrat cible, pas un état des fonctionnalités réalisées. Toute feature doit être accompagnée de critères d'acceptation et d'une matrice de compatibilité validée.

## Ordre de lecture

1. architecture.md : frontières et orchestration.
2. hardware-access.md : sources, privilèges et limites.
3. data-model.md : unités, provenance et snapshots.
4. cli.md : interface utilisateur et automatisation.
5. recommendations.md : règles et garde-fous.
6. roadmap.md : features, dépendances et acceptance.
7. ../CONTRIBUTING.md et ../SECURITY.md : qualité et sécurité.

## Décisions initiales

- D001 : CLI en Go, version cible 1.27 sous réserve de disponibilité vérifiée.
- D002 : Windows 11 x64 comme première cible d'acquisition.
- D003 : collecte en lecture seule et aucune modification automatique du BIOS/SPD.
- D004 : données partielles acceptées, absence explicite et provenance obligatoire.
- D005 : CSV long et snapshots JSON versionnés.
- D006 : aucun pilote vulnérable ni désactivation de HVCI/Secure Boot.
- D007 : SDK CPUID possible uniquement après audit technique et juridique ; aucune dépendance décidée.
- D008 : recommandations manuelles expérimentales conditionnées à des règles validées par plateforme.

## Décisions ouvertes

Licence, backend SPD/runtime, bibliothèque CLI, stratégie de signature des composants natifs, règles d'optimisation par plateforme, formats des résultats de tests.

Une décision importante devra produire un ADR : contexte, alternatives, preuves, décision, conséquences, statut et date. Créer docs/adr/ lors de la première décision complémentaire.

## Références

Les références techniques sont centralisées dans hardware-access.md. Consigner les versions de spécifications effectivement utilisées, leurs conditions de redistribution et la date de consultation lors de l'implémentation.
