package pastebin

type PastebinOpt func(*Pastebin) error

func PastebinLoginOpt(userKey string) PastebinOpt {
	return func(p *Pastebin) error {
		p.userKey = userKey
		return nil
	}
}

func PastebinFolderOpt(folderKey string) PastebinOpt {
	return func(p *Pastebin) error {
		p.folder = folderKey
		return nil
	}
}
