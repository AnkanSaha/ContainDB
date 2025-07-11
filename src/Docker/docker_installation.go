package Docker

import (
	"fmt"
	"os"
	"os/exec"
)

func IsDockerInstalled() bool {
	cmd := exec.Command("docker", "--version")
	err := cmd.Run()
	return err == nil
}

func InstallDocker() error {
	fmt.Println("Docker not found. Installing Docker...")
	commands := []string{
		"sudo apt-get update",
		"sudo apt-get install -y ca-certificates curl",
		"sudo install -m 0755 -d /etc/apt/keyrings",
		"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
		"sudo chmod a+r /etc/apt/keyrings/docker.asc",
		`echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
		$(. /etc/os-release && echo ${UBUNTU_CODENAME:-$VERSION_CODENAME}) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null`,
		"sudo apt-get update",
		"sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
		"sudo usermod -aG docker $USER",
		"sudo apt install -y docker-compose-plugin",
	}

	for index, c := range commands {
		fmt.Println("Running command", index+1, "of", len(commands), ":", c)
		cmd := exec.Command("bash", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running command:", c)
		}
	}
	return nil
}

func UninstallDocker() error {
	fmt.Println("Uninstalling Docker...")
	commands := []string{
		"sudo apt-get purge -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
		"sudo rm -rf /var/lib/docker",
		"sudo rm -rf /var/lib/containerd",
		"sudo rm /etc/apt/sources.list.d/docker.list",
		"sudo rm /etc/apt/keyrings/docker.asc",
		"sudo apt-get autoremove -y",
	}

	for index, c := range commands {
		fmt.Println("Running command", index+1, "of", len(commands), ":", c)
		cmd := exec.Command("bash", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running command:", c)
		}
	}
	return nil
}
