# Politique de recommandations

## Objectif

Proposer des candidats BIOS explicables. Pas de réglage universel garanti, pas d'écriture matérielle, pas de promesse de gain basée uniquement sur CL.

## Entrées nécessaires

CPU exact, carte mère/révision, BIOS, nombre et emplacement des DIMM, références, profils SPD et paramètres actifs disponibles. Pour des règles manuelles : identification mémoire suffisamment fiable, paramètres runtime concernés, enveloppe de tension documentée pour la plateforme et historique de validation.

Une QVL est une preuve de validation dans les conditions du fabricant, pas une garantie universelle. Conserver URL, date, version et conditions si utilisée. Ne pas scraper ou prétendre posséder toutes les QVL dans la première version.

## Catégories

- manufacturer_profile : paramètres lus dans un profil constructeur valide.
- experimental_candidate : hypothèse issue d'une règle versionnée applicable à la plateforme.
- tested_configuration : configuration associée à des tests précis sur cette machine.

Résultats : candidate, insufficient_data, unsupported, rejected. La réussite de tests ne garantit pas toute charge future ; préciser outils, versions, durée, erreurs, conditions et portée.

## Sélection de profil

Valider chaque profil et sa cohérence. Comparer les configurations de tous les modules, pas simplement leurs noms de profil. Distinguer kit unique et mélange de kits. Ne pas composer automatiquement les plus petits timings de profils différents.

Identifier les manques de compatibilité CPU/BIOS. L'absence dans une QVL ne prouve pas l'incompatibilité. Comparer débit, timings et tensions sans classer tous les kits par CL seul. Si les paramètres actifs correspondent, afficher une correspondance possible, pas une activation certaine.

## Réglages manuels

Désactivés par défaut jusqu'à existence d'une règle vérifiée pour la plateforme. Aucune règle générique pour tension SoC, VDD, VDDQ ou contrôleur mémoire. Ne pas recommander un seuil dit sûr sans référence autorisée et contexte exact.

Toute règle contient ID/version, sélecteur matériel, préconditions, sources, paramètres couplés, limites, justification et méthode de validation. Un manque de données bloque la proposition concernée.

## Validation

Conserver une référence stable et la procédure de récupération BIOS du constructeur. Changer une variable ou un ensemble couplé documenté à la fois. Après changement : nouvelle capture, tests mémoire avec outils externes choisis, redémarrages et charges représentatives. Mesurer le gain avec plusieurs répétitions et conserver les conditions. Un POST réussi ou un meilleur CL calculé n'est pas un verdict.

## Sortie structurée

rule_id, rule_version, category, status, prerequisites, evidence, proposed_settings, missing_data, warnings, validation_plan. Les tensions et timings proposés portent leur origine. Les recommandations ne sont pas des valeurs actuelles et ne figurent jamais dans l'export current.
