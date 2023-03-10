package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/halilylm/micro/foundation/vault"
	"log"
	"os"
	"time"
)

const credentialsFileName = "/vault/credentials.json"

// VaultInit sets up a newly provisioned vault instance.
func VaultInit(vaultConfig vault.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vaultSrv, err := vault.New(vault.Config{
		Address:   vaultConfig.Address,
		MountPath: vaultConfig.MountPath,
	})
	if err != nil {
		return fmt.Errorf("constructing vault: %w", err)
	}

	initResponse, err := checkIfCredFileExists()
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			log.Println("credential file doesn't exists, initializing vault")

			initResponse, err = vaultSrv.SystemInit(ctx, 1, 1)
			if err != nil {
				if errors.Is(err, vault.ErrAlreadyInitialized) {
					return fmt.Errorf("vault is already initialized but we do not have credentials file")
				}
				return fmt.Errorf("unable to initialize Vault instance: %w", err)
			}

			b, err := json.Marshal(initResponse)
			if err != nil {
				return errors.New("unable to marshal")
			}

			if err := os.WriteFile(credentialsFileName, b, 0644); err != nil {
				return fmt.Errorf("unable to write %s file: %w", credentialsFileName, err)
			}
		default:
			return fmt.Errorf("unable to read credentials file: %w", err)
		}
	}
	log.Printf("rootToken: %s", initResponse.RootToken)

	log.Println("unsealing vault")
	err = vaultSrv.Unseal(ctx, initResponse.KeysB64[0])
	if err != nil {
		if errors.Is(err, vault.ErrBadRequest) {
			return fmt.Errorf("vault is not initialized. Check for old credentials file: %s", credentialsFileName)
		}
		return fmt.Errorf("error unsealing vault: %w", err)
	}

	log.Println("mounting path in vault")

	vaultSrv.SetToken(initResponse.RootToken)
	if err := vaultSrv.Mount(ctx); err != nil {
		if errors.Is(err, vault.ErrPathInUse) {
			return fmt.Errorf("unable to mount the path: %w", err)
		}
		return fmt.Errorf("error unsealing vault: %w", err)
	}

	log.Println("creating sales-api policy")

	err = vaultSrv.CreatePolicy(ctx, "sales-api", "secret/data/*", []string{"read", "create", "update"})
	if err != nil {
		return fmt.Errorf("unable to create policy: %w", err)
	}

	log.Println("Generating sales-api token: %s", vaultConfig.Token)

	err = vaultSrv.CheckToken(ctx, vaultConfig.Token)
	if err == nil {
		log.Println("token already exists: ", vaultConfig.Token)
		return nil
	}

	err = vaultSrv.CreateToken(ctx, vaultConfig.Token, []string{"sales-api"}, "Sales API")
	if err != nil {
		return fmt.Errorf("unable to create token: %w", err)
	}

	return nil
}

func checkIfCredFileExists() (vault.SystemInitResponse, error) {
	if _, err := os.Stat(credentialsFileName); err != nil {
		return vault.SystemInitResponse{}, err
	}

	data, err := os.ReadFile(credentialsFileName)
	if err != nil {
		return vault.SystemInitResponse{}, fmt.Errorf("reading %s file: %w", credentialsFileName, err)
	}

	var response vault.SystemInitResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return vault.SystemInitResponse{}, fmt.Errorf("unmarshalling json: %w", err)
	}
	return response, nil
}