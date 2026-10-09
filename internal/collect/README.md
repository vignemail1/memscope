# internal/collect

Orchestration indépendante de Windows, sans accès matériel propre ni élévation automatique. Les interfaces de provider.go restent inchangées.

## Utilisation

```go
collector, err := collect.New(providers, collect.Config{
    SnapshotID: captureID,
    ToolVersion: version,
    SchemaVersion: "1",
    Timeout: 10 * time.Second,
})
if err != nil {
    return err
}
snapshot, err := collector.Collect(ctx)
// Toujours examiner snapshot ET err : les résultats partiels sont conservés.
```

L'identifiant est fourni par l'appelant. Le package ne génère pas d'UUID et n'interprète pas la compatibilité des versions persistées.

## Contrat

- Ordre séquentiel : inventory, SPD, runtime, pour éviter les accès bus concurrents dans cette orchestration.
- Un provider absent ajoute une capacité indisponible mais n'est pas une erreur.
- Une erreur de probe est conservée et empêche l'acquisition du provider concerné.
- Les capacités sont informatives : aucun nom de capacité spécifique ne pilote implicitement une opération ; le provider contrôle son acquisition.
- Une erreur d'acquisition ne supprime pas les observations, profils ou diagnostics partiels.
- Les erreurs sont contextualisées et agrégées via errors.Join, avec errors.Is utilisable.
- La validation finale peut échouer ; le snapshot brut reste retourné pour diagnostic et ne doit pas être exporté comme une capture valide.
- Les observations contradictoires ne sont ni fusionnées ni corrigées silencieusement.

Les SPD/runtime doivent référencer des appareils fournis par l'inventaire ; leurs interfaces ne retournent pas de nouveaux appareils. Si l'inventaire est incomplet, les références non résolues entraînent une erreur de validation plutôt qu'une identité inventée.

## Limites

Timeout coopératif : un appel natif qui ignore le contexte peut bloquer. Un futur backend natif devra fournir des timeouts propres ou une isolation de processus. Aucune goroutine abandonnée pour simuler une interruption.

Un provider absent doit être une interface nil, pas une interface contenant un pointeur typé nil. Les providers ne doivent pas muter les données retournées. Collector ne garantit pas leur thread-safety ; ne pas lancer plusieurs collectes concurrentes sur les mêmes providers sans protection backend.

Ne pas réutiliser le même SnapshotID pour plusieurs captures. Ce package ne vérifie pas l'unicité globale. Les diagnostics techniques doivent être assainis par les providers avant de contenir des données personnelles.

## Tests

Tests ajoutés : configuration, providers absents, résultats partiels, probe en échec, annulation, deadline coopérative, snapshot invalide et entrées nil.

```sh
go test ./internal/collect
go vet ./internal/collect
go test ./...
```

Ces contrôles n'ont pas été exécutés lors de la préparation du commit. Les collecteurs matériels et l'intégration CLI ne sont pas implémentés par ce changement.
