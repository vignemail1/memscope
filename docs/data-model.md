# Modèle de données et export

## Snapshot JSON v1

Champs cibles : schema_version, tool_version, snapshot_id, collection_started_at, collection_finished_at, platform, devices, observations, profiles, diagnostics, capabilities. Identifiants locaux stables dans le snapshot, pas nécessairement entre changements matériels.

Une observation contient device_id, scope, parameter, value nullable, unit, source, source_version, status, captured_at, derivation facultative et diagnostics.

Sources : smbios, wmi, cpuid, spd, controller, telemetry, calculation, import, user. Statuts : observed, reported, derived, unavailable, invalid. La confiance est distincte du statut ; une valeur reported n'est pas forcément fausse.

Les profils appartiennent aux DIMM. Les timings actifs appartiennent au contrôleur/canal concerné ; ne pas les dupliquer sur toutes les DIMM comme des mesures indépendantes. Conserver valeurs contradictoires et raison du choix d'affichage.

## Unités et calculs

Débit DDR : MT/s ; horloge mémoire : MHz ; timings : cycles ; durées : ns ; tension : V ; capacité : bytes. Utiliser des identifiants canoniques et une précision conservée.

Pour une horloge correspondant bien à la convention DDR utilisée : délai_ns = cycles * 2000 / débit_MT_s. Exemple : 5600 CL46 donne 16,428571 ns. Calculer seulement avec entrées valides, positives et de provenance connue. Ne pas assimiler ce délai à la latence complète CPU-vers-RAM. Ne pas utiliser un ratio contrôleur comme horloge DRAM.

Une tension de profil reste reported dans le périmètre profile. Une consigne runtime et une tension mesurée sont des paramètres distincts.

## CSV current v1

Colonnes fixes : schema_version,snapshot_id,captured_at,scope,device_id,parameter,value,unit,source,status

Format long : une ligne par observation, ordonnée par scope, device_id, parameter et source. Champs inconnus vides, status explicite. UTF-8 sans BOM par défaut, virgule, point décimal, fins de ligne CRLF. Échappement via encoding/csv. Les timestamps utilisent RFC3339 avec fuseau.

## CSV profiles v1

Colonnes fixes : schema_version,snapshot_id,device_id,profile_id,profile_type,parameter,value,unit,source,status

Une ligne par paramètre de profil. profile_type distingue JEDEC/XMP/EXPO et variantes décodées sans perdre leur version dans le modèle JSON. Ne pas fabriquer une configuration complète à partir de contraintes partielles.

## Protection et confidentialité

Exclure par défaut numéros de série, UUID machine et identifiants personnels des exports. Leur présence exige --include-sensitive. Les journaux n'incluent jamais le dump brut sans consentement. Les dumps bruts sont hors snapshots usuels.

Neutraliser les cellules textuelles externes pouvant être interprétées comme formules par un tableur, y compris après espaces/contrôles initiaux ; conserver l'original dans le modèle interne. Les valeurs numériques validées sont sérialisées comme nombres. Documenter l'altération de présentation et tester les caractères =,+,-,@, tabulation et CR.

## Versionnement

Tout changement cassant augmente schema_version. Le lecteur refuse les versions majeures non prises en charge. Définir une politique explicite pour champs supplémentaires et migrations ; pas de conversion silencieuse perdant la provenance. Fixtures de référence pour chaque version.
