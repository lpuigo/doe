package manager

import "io"

func (m Manager) ClientsArchiveName() string {
	return m.Clients.ArchiveName()
}

func (m Manager) CreateClientsArchive(writer io.Writer) error {
	return m.Clients.CreateArchive(writer)
}
