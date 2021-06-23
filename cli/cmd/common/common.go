package common

import (
	"fmt"
	root "github.com/timescale/tobs/cli/cmd"
	"github.com/timescale/tobs/cli/pkg/helm"
	"github.com/timescale/tobs/cli/pkg/pgconn"
)

const (
	LISTEN_PORT_GRAFANA    = 8080
	FORWARD_PORT_GRAFANA   = 3000
	LISTEN_PORT_PROM       = 9090
	FORWARD_PORT_PROM      = 9090
	LISTEN_PORT_PROMLENS   = 8081
	FORWARD_PORT_PROMLENS  = 8080
	LISTEN_PORT_PROMSCALE  = 9201
	FORWARD_PORT_PROMSCALE = 9201
	LISTEN_PORT_TSDB       = 5432
	FORWARD_PORT_TSDB      = 5432
)

func FormDBDetails(user, dbName string) (pgconn.DBDetails, error) {
	secretKey, user, err := GetDBSecretKeyAndDBUser(root.HelmReleaseName, user)
	if err != nil {
		return pgconn.DBDetails{}, fmt.Errorf("could not get DB secret key from helm release: %w", err)
	}

	d := pgconn.DBDetails{
		Namespace: root.Namespace,
		Name:      root.HelmReleaseName,
		DBName:    dbName,
		User:      user,
		SecretKey: secretKey,
		Remote:    FORWARD_PORT_TSDB,
	}

	return d, nil
}

func GetDBSecretKeyAndDBUser(releaseName, dbUser string) (string, string, error) {
	var userName string
	helmClient := helm.NewClient(root.Namespace)
	e, err := helmClient.ExportValuesFieldFromRelease(releaseName, []string{"timescaledb-single", "enabled"})
	if err != nil {
		return "", "", err
	}
	enableTimescaleDB, ok := e.(bool)
	if !ok {
		return "", "", fmt.Errorf("enable Backup was not a bool")
	}

	if !enableTimescaleDB {
		dbURI, err := helmClient.ExportValuesFieldFromRelease(releaseName, []string{"timescaledbExternal", "db_uri"})
		if err != nil {
			return "", "", err
		}

		uriDetails, err := pgconn.ParseDBURI(fmt.Sprint(dbURI))
		if err != nil {
			return "", "", err
		}
		userName = uriDetails.ConnConfig.User
		return "PATRONI_SUPERUSER_PASSWORD", userName, nil
	}

	data, err := helmClient.ExportValuesFieldFromRelease(releaseName, []string{"timescaledb-single", "patroni", "postgresql", "authentication", "superuser", "username"})
	if err != nil {
		return "", "", err
	}

	userName = fmt.Sprint(data)

	// fetch the superUser from deployment
	// if dbUser is not empty && superUser != user then return provided user as secretKey & user
	// else send the default secretKey & superUser fetched
	if dbUser != "" && dbUser != userName {
		return dbUser, dbUser, nil
	}

	// As the user isn't provided
	// use default super user from helm release
	// the default super-user password is mapped to "PATRONI_SUPERUSER_PASSWORD" secret key
	return "PATRONI_SUPERUSER_PASSWORD", fmt.Sprint(userName), nil
}