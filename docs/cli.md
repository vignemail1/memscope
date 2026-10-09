# Contrat CLI proposé

Toutes les commandes sont spécifiées, pas encore implémentées.

## Commandes

| Commande | Contrat |
|---|---|
| inspect | Résumé matériel, provenance et limites |
| memory current | Valeurs actives seulement ; firmware déclaré identifié comme tel |
| memory profiles | Profils/contraintes SPD par module, sans présumer activation |
| recommend | Candidats BIOS et données manquantes ; aucune modification |
| doctor | Capacités, privilèges, composants et raisons des indisponibilités |
| export | CSV current ou profiles depuis acquisition ou snapshot |
| snapshot | Capture JSON versionnée |
| compare | Différences entre deux snapshots compatibles |

## Options

Options globales : --config, --timeout, --verbose, --no-color. Ajouter --version et --help. Les options devront être acceptées avant ou après les sous-commandes.

inspect, memory current, memory profiles et recommend : --format text|json et --input snapshot.json pour l'analyse hors ligne.

doctor : --format text|json.

export : --scope current|profiles, --format csv, --output PATH obligatoire, --input PATH facultatif, --force, --include-sensitive.

snapshot : --output PATH obligatoire, --force, --include-sensitive.

compare : BEFORE AFTER, --format text|json.

Données stdout, diagnostics stderr. Aucun code couleur ni barre de progression dans les sorties structurées. Aucune invite interactive par défaut. Le mode JSON doit garder une enveloppe stable même si la collecte est partielle.

## Configuration

Priorité : options explicites > variables MEMSCOPE_* > fichier explicitement demandé > valeurs par défaut. Prévoir MEMSCOPE_TIMEOUT, MEMSCOPE_LOG_LEVEL et MEMSCOPE_CONFIG. Ne pas chercher un fichier arbitraire dans le répertoire courant. Documenter toutes les clés au moment de leur implémentation.

## Codes de sortie

- 0 : opération réalisée, y compris résultat partiel documenté.
- 1 : erreur opérationnelle empêchant le résultat demandé.
- 2 : arguments ou configuration invalides.
- 3 : capacité indispensable indisponible.
- 4 : permission insuffisante pour une capacité indispensable.
- 5 : snapshot invalide ou version incompatible.

La complétude figure dans la sortie ; un CSV partiel avec champs absents n'est pas automatiquement un échec. Une recommandation impossible indique status=insufficient_data et liste les manques.

## Fichiers et comparaison

Pas d'écrasement sans --force. Écriture temporaire dans le même répertoire puis publication selon garanties de la plateforme. Pas d'écriture partielle du fichier cible. Comparer unités, périmètre et identité matérielle avant les valeurs ; signaler les changements de matériel et ne pas comparer implicitement deux DIMM différents.
