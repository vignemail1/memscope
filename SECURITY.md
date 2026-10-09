# Sécurité

## Engagements de conception

- Aucun changement automatique de BIOS, timing, tension ou SPD.
- Aucun contournement de HVCI, Secure Boot ou de la liste de pilotes bloqués.
- Aucun pilote vulnérable, installation silencieuse ou élévation automatique.
- Accès privilégié minimal, opérations en liste blanche et validation stricte des entrées.
- Aucune télémétrie réseau par défaut.

## Menaces à traiter

Firmware/SPD malformé, fichiers importés hostiles, allocations excessives, chemins d'export, injection CSV, DLL détournée, IPC non autorisé, accès bus concurrent, confusion profil/mesure et recommandation dangereuse.

Les composants natifs doivent être chargés depuis un emplacement maîtrisé avec vérification de provenance. Les fixtures et dumps peuvent contenir des identifiants personnels ; anonymisation avant publication.

## Signalement

Ne pas publier d'exploit, dump sensible ou secret dans une issue publique. Aucun canal privé spécifique n'est encore configuré : le propriétaire doit activer et annoncer un moyen de signalement privé avant la première release. Ce document ne prétend pas que GitHub Private Vulnerability Reporting est déjà activé.

## Versions prises en charge

Aucune release à ce stade. La politique de support et les procédures de signature/mise à jour seront définies avant distribution d'un binaire ou pilote.
