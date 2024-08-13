package atmoz_sftp_test

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/sftp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSFTPServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SFTP Server Suite")
}

var _ = Describe("SFTP Server", func() {
	var (
		ctx           context.Context
		sftpContainer testcontainers.Container
		sshClient     *ssh.Client
		sftpClient    *sftp.Client
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		sftpContainer, err = startSFTPContainer(ctx)
		Expect(err).NotTo(HaveOccurred(), "Failed to start SFTP container")

		sshClient, err = connectToSFTP(ctx, sftpContainer)
		Expect(err).NotTo(HaveOccurred(), "Failed to connect to SFTP server")

		sftpClient, err = sftp.NewClient(sshClient)
		Expect(err).NotTo(HaveOccurred(), "Failed to create SFTP sshClient")
	})

	AfterEach(func() {
		Expect(sftpClient.Close()).To(Succeed(), "Failed to close SFTP sshClient")
		Expect(sshClient.Close()).To(Succeed(), "Failed to close SSH sshClient")
		Expect(sftpContainer.Terminate(ctx)).To(Succeed(), "Failed to terminate SFTP container")
	})

	Context("Uploading and verifying a file", func() {
		It("should upload and verify the file content", func() {
			err := uploadAndVerifyFile(sftpClient)
			Expect(err).NotTo(HaveOccurred(), "Failed to upload and verify file on SFTP server")
		})
	})
})

func startSFTPContainer(ctx context.Context) (testcontainers.Container, error) {
	containerRequest := testcontainers.ContainerRequest{
		Image:        "atmoz/sftp",
		ExposedPorts: []string{"22/tcp"},
		Env: map[string]string{
			"SFTP_USERS": "foo:pass:::upload",
		},
		WaitingFor: wait.ForListeningPort("22/tcp").
			WithStartupTimeout(10 * time.Second),
	}

	sftpContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	return sftpContainer, nil
}

func connectToSFTP(ctx context.Context, container testcontainers.Container) (*ssh.Client, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "22")
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: "foo",
		Auth: []ssh.AuthMethod{
			ssh.Password("pass"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := net.JoinHostPort(host, port.Port())
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func uploadAndVerifyFile(sftpClient *sftp.Client) error {
	tmpFile, err := os.CreateTemp("", "example")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("Hello, SFTP!")
	if _, err = tmpFile.Write(content); err != nil {
		return err
	}

	dstPath := "/upload/testfile.txt"
	dstFile, err := sftpClient.Create(dstPath)
	if err != nil {
		return err
	}

	if _, err = dstFile.Write(content); err != nil {
		return err
	}
	dstFile.Close()

	srcFile, err := sftpClient.Open(dstPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	uploadedContent, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	if string(content) != string(uploadedContent) {
		return fmt.Errorf("content mismatch: expected %s, got %s", content, uploadedContent)
	}

	return nil
}
