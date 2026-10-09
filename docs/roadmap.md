# Features et roadmap

Toutes les features ci-dessous sont à réaliser. Priorités P0 fondations, P1 produit mémoire, P2 expérimentation, P3 extension.

## M0 — Décisions préalables

- F001 P0 : vérifier Go 1.27 et figer toolchain/outils. Acceptance : disponibilité et versions consignées ; aucune préversion implicite.
- F002 P0 : choisir licence et politique de dépendances. Acceptance : audit des composants et droits de redistribution.
- F003 P0 : étude backend matériel. Acceptance : comparaison CPUID/indépendant/import avec licences, HVCI, données et matériel testable.

## M1 — Socle CLI

- F010 P0 : modèle versionné, provenance et diagnostics. Acceptance : fixtures partielles/contradictoires sans valeurs inventées.
- F011 P0 : commandes, aide, configuration et codes de sortie. Acceptance : tests de priorité, stdout/stderr et non-interactivité.
- F012 P0 : inventaire SMBIOS Windows. Acceptance : longueurs/sentinelles testées, collecte utilisateur standard et résultats comparés au firmware.
- F013 P0 : doctor. Acceptance : distingue absence, accès refusé, plateforme non supportée et erreur de backend.
- F014 P0 : JSON/CSV et compare. Acceptance : schemas stables, tests golden, confidentialité, injection CSV et écriture sans écrasement.
- F015 P0 : CI et build Go. Acceptance : format, vet, lint, tests et analyses de vulnérabilités passent ; artefacts identifiables.

M1 n'est pas un équivalent CPU-Z et ne produit pas de recommandations SPD natives.

## M2 — SPD et état actif

- F020 P1 : transport SPD validé. Acceptance : limites/timeouts, concurrence, aucune écriture persistante et matrice chipset.
- F021 P1 : décodage DDR4. Acceptance : spécification versionnée, fixtures autorisées, CRC/révisions et entrées corrompues.
- F022 P1 : décodage DDR5 et SPD hub. Acceptance : tests pagination, révisions et données partielles.
- F023 P1 : extensions XMP/EXPO. Acceptance : profil invalidé si incohérent ; paramètres/tensions/provenance exportables.
- F024 P1 : paramètres runtime. Acceptance : support CPU/chipset explicite, distinction tRCD lecture/écriture si applicable et comparaison à outil de référence.
- F025 P1 : tensions. Acceptance : mesure, consigne et profil séparés ; unavailable si non pris en charge.
- F026 P1 : tableau de compatibilité réel. Acceptance : OS/build, CPU, carte mère, BIOS, modules, backend, HVCI et capacités documentés.

## M3 — Recommandations constructeur

- F030 P1 : moteur de règles versionné. Acceptance : déterministe sur snapshots identiques, explications et préconditions.
- F031 P1 : candidats XMP/EXPO communs. Acceptance : modules mixtes, profils absents et réglages déjà correspondants testés.
- F032 P1 : plan de validation. Acceptance : pas de garantie de stabilité ni tension générique ; données manquantes explicites.

## M4 — Validation et optimisation expérimentale

- F040 P2 : attacher résultats de tests externes. Acceptance : versions, durées, erreurs et configuration exacte conservées.
- F041 P2 : comparer benchmarks répétés. Acceptance : statistiques et conditions, pas de conclusion sur un unique score.
- F042 P2 : règles manuelles par plateforme. Acceptance : sources, limites, tests et désactivation hors sélecteur ; review sécurité obligatoire.

## M5 — Couverture matérielle élargie

- F050 P3 : CPU caches, instructions et topologie détaillée.
- F051 P3 : timings secondaires/tertiaires supplémentaires.
- F052 P3 : backends et architectures supplémentaires après validation.

## Definition of Done

Code revu, format/lint/vet/tests propres, risques et licence évalués, fixtures anonymisées, documentation à jour, comportement dégradé testé et feature annoncée uniquement sur la matrice validée. Les tests logiciels ne remplacent pas la validation physique des collecteurs.

Pas de publication de pilote, release binaire ou activation de règle expérimentale sans revue explicite. Transformer les Fxxx en issues lors du démarrage de chaque jalon ; ce commit ne crée pas d'issues GitHub.
