# Méthodes d'accès aux informations

## Sources et limites

| Source | Informations candidates | Limites |
|---|---|---|
| GetSystemFirmwareTable avec RSMB | SMBIOS brut : BIOS, système, carte mère, processeur, appareils mémoire | Déclarations firmware, pas lecture SPD ni timings actifs universels |
| WMI/CIM | Win32_PhysicalMemory, Win32_Processor, Win32_BaseBoard, Win32_BIOS | Schéma dépendant du firmware ; pas de table complète XMP/EXPO |
| CPUID CPU | Identité, caractéristiques, caches/topologie selon architecture | Ne lit pas les profils SPD ; ne suffit pas aux timings mémoire actifs |
| SPD via transport matériel validé | Identité DIMM, organisation, contraintes JEDEC et extensions XMP/EXPO | Dépend du chipset, accès bus, module et spécification |
| Registres contrôleur/SDK | Paramètres mémoire actifs disponibles | Spécifiques plateforme ; privilèges et documentation nécessaires |
| Télémétrie carte mère/PMIC/SDK | Tensions réellement observées lorsque prises en charge | Ne pas confondre consigne, profil et mesure |

## Inventaire initial

Privilégier GetSystemFirmwareTable pour l'accès SMBIOS brut ; évaluer WMI comme complément. Parser les longueurs, chaînes, sentinelles et champs dépendant de la version SMBIOS. Prendre en charge les valeurs étendues prévues par la spécification utilisée. Ne pas interpréter une absence comme zéro. Préserver unités et données brutes utiles.

Ne pas assimiler ConfiguredClockSpeed à une mesure temps réel. Documenter l'interprétation firmware et les unités ; ne jamais appliquer automatiquement un facteur deux sans base explicite.

## SPD

Le SPD n'est pas le SMBIOS. Pour DDR4 et DDR5, transports et organisation diffèrent ; DDR5 utilise un SPD hub. Pas de scan agressif ni de lecture d'adresses arbitraires. Établir une matrice chipset/contrôleur, gérer pagination, accès concurrents, timeouts et tailles limites.

Distinguer programmation du SPD, interdite, et transactions de sélection de page éventuellement nécessaires à la lecture : celles-ci doivent être documentées, minimisées et ne modifier aucune configuration persistante.

Le décodeur valide type/version, tailles, CRC lorsque défini et plages avant d'exposer un profil. Les paramètres JEDEC ne constituent pas nécessairement une suite exhaustive de profils discrets : respecter les contraintes de temps, masques de CL et règles de conversion de la spécification. Ne pas inventer un tableau complet.

## Paramètres actifs

CL, tRCD, tRP, tRAS, command rate, timings secondaires et ratios ne viennent pas simplement du profil sélectionné. Lire une source runtime documentée. Certaines plateformes séparent tRCD lecture/écriture : préserver les champs distincts. L'identification d'un profil actif reste une correspondance candidate, pas une preuve, si aucun indicateur explicite n'existe.

## Choix du backend

Option A : évaluer le SDK CPUID (couverture, distribution, coûts, support, pilotes et licence). CPU-Z freeware n'implique pas SDK libre ni API publique réutilisable.

Option B : backend indépendant avec documentation autorisée et composants audités. Un projet existant ne peut être réutilisé qu'après audit licence, sécurité, signature et maintenance.

Option C : imports documentés pour développement hors matériel ; ne pas les présenter comme une acquisition native.

Livrable préalable : matrice comparant couverture SPD/runtime/tensions, plateformes, privilèges, HVCI, licence, coût de maintenance et tests réalisables.

## Sécurité

Aucun contournement de HVCI, Secure Boot ou blocklist. Aucun ancien pilote vulnérable pour gagner de la couverture. Pas d'accès matériel arbitraire, d'installation silencieuse ou d'élévation automatique. Une source indisponible produit un diagnostic exploitable.

## Références de travail

- Microsoft GetSystemFirmwareTable : https://learn.microsoft.com/en-us/windows/win32/api/sysinfoapi/nf-sysinfoapi-getsystemfirmwaretable
- Microsoft Win32_PhysicalMemory : https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-physicalmemory
- Microsoft driver block rules : https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/design/microsoft-recommended-driver-block-rules
- DMTF SMBIOS : https://www.dmtf.org/standards/smbios
- CPUID : https://www.cpuid.com/
- Intel XMP : https://www.intel.com/content/www/us/en/gaming/extreme-memory-profile-xmp.html

Les offsets SPD/XMP/EXPO et registres runtime doivent être fondés sur une spécification versionnée accessible légalement. Ne pas les déduire de snippets ou de ce document.
