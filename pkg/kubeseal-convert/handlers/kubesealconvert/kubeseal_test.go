package kubesealconvert

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
)

func Test_buildDataBytes(t *testing.T) {
	tests := []struct {
		name string
		sv   domain.SecretValues
		want map[string][]byte
	}{
		{
			name: "single string value",
			sv: domain.SecretValues{
				Data: map[string]interface{}{"password": "s3cr3t"},
			},
			want: map[string][]byte{"password": []byte("s3cr3t")},
		},
		{
			name: "multiple string values",
			sv: domain.SecretValues{
				Data: map[string]interface{}{
					"username": "admin",
					"password": "s3cr3t",
				},
			},
			want: map[string][]byte{
				"username": []byte("admin"),
				"password": []byte("s3cr3t"),
			},
		},
		{
			name: "empty data",
			sv:   domain.SecretValues{Data: map[string]interface{}{}},
			want: map[string][]byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDataBytes(tt.sv)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_buildKubeSecret(t *testing.T) {
	sv := domain.SecretValues{
		Name:      "my-secret",
		Namespace: "production",
		Data:      map[string]interface{}{"token": "abc123"},
		Annotations: map[string]string{
			"note": "generated",
		},
		Labels: map[string]string{
			"env": "prod",
		},
	}

	got := buildKubeSecret(sv)

	assert.Equal(t, "my-secret", got.Name)
	assert.Equal(t, "production", got.Namespace)
	assert.Equal(t, "v1", got.APIVersion)
	assert.Equal(t, "Secret", got.Kind)
	assert.Equal(t, "Opaque", string(got.Type))
	assert.Equal(t, map[string][]byte{"token": []byte("abc123")}, got.Data)
	assert.Equal(t, map[string]string{"note": "generated"}, got.Annotations)
	assert.Equal(t, map[string]string{"env": "prod"}, got.Labels)
}

func Test_runCommandWithInput(t *testing.T) {
	echoBin, err := exec.LookPath("echo")
	require.NoError(t, err, "echo binary must be available")

	tests := []struct {
		name    string
		cmdPath string
		args    []string
		input   string
		wantOut string
		wantErr bool
	}{
		{
			name:    "echo passes stdin through args (no stdin capture)",
			cmdPath: echoBin,
			args:    []string{"hello"},
			input:   "ignored-stdin",
			wantOut: "hello\n",
		},
		{
			name:    "invalid binary path returns error",
			cmdPath: "/nonexistent/binary",
			args:    []string{},
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := runCommandWithInput(tt.cmdPath, tt.args, tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOut, out)
			}
		})
	}
}
