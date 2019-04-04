# Filemux

## Construction

```bash
# GO111MODULE=on go mod vendor # Rebuild des vendors.

CGO_ENABLED=0 GOOS=linux go build -o dist/linux/filemux # Build pour Linux.
CGO_ENABLED=0 GOOS=windows go build -o dist/windows/filemux.exe # Build pour Windows.
```

## Lancer les tests

```bash
go test ./...
```

## Fichier de configuration

Des exemples sont disponibles [ici](samples/etc/filemux).

- La paire key/value `out` représente une liste d'objets où chaque objet représente un ticker qui (périodiquement configuré avec `interval`) vérifie la présence de fichiers dans le glob `${src}/*`. Chaque fois que le tick se déclenche, il vérifie s'il existe des fichiers disponibles pour l'archivage (chaque nom de fichier dans l'archive peut être préfixé via `archive_inner_dirname`) dans `${dst}/${nom_base_archive}.tar`. Les fichiers écrits avec succès dans l'archive sont ensuite supprimés de `src`. La commande `exec` est finalement exécutée.
- La paire key/value `in` représente une liste d'objets où chaque objet représente un notificateur qui écoute les événements de type création sur `src`. Chaque fois qu'un nouveau fichier est créé, une routine est lancée pour vérifier que l'écriture de ce fichier est terminée. L'archive tar est finalement extraite à `dst`.

Tous les templates Go sont injectés avec les [helpers sprig](http://masterminds.github.io/sprig).

## Validation d'un fichier de configuration

```bash
filemux validate -c filemux.yml
```

## Lancement

```bash
# filemux help # Affiche l'aide globale.
# filemux help run # Affiche l'aide pour la sous-commande run.

filemux run -c filemux.yml
```

## Lancement via un service `initd`

Voir [ici](samples/etc/init.d/filemux) pour le *service unit file*.

```bash
chkconfig filemux on
service filemux start
```

## Lancement via un service `systemd`

Voir [ici](samples/etc/systemd/system/filemux.service) pour le *service unit file*.

```bash
systemctl enable filemux
systemctl start filemux
```

## Lancement via un service Win32 

1. [Télécharger le *Non-Sucking Service Manager*](https://nssm.cc/download).
2.
```cmd
nssm install filemux C:\filemux\filemux.exe run -c C:\filemux\filemux.yml
nssm set filemux AppStdout C:\filemux\log\filemux.log

REM "net stop filemux" : 10000ms est le temps laissé au service pour un arrêt en douceur avant l'appel de TerminateProcess().
nssm set filemux AppStopMethodConsole 10000 

REM Le service sera démarré avec l'utilisateur filemux
nssm get filemux ObjectName filemux <password>

net start filemux
```
