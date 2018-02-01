package dao

import "database/sql"

type DAO struct {
	Db *sql.DB
}

func (d *DAO) Comments(threadHash string) ([]*Comment, error) {
	rows, err := d.Db.Query(`
    		select  coalesce(
    			(select name from addresses where address=sender), sender
    		), comment, received,
    		sender, body, pending, subject, extended, id, thread_hash
            from channels
            where thread_hash = ? and received > unix_timestamp() - 6*28*24*60*60
            order by received asc
            limit 1000
    	`, threadHash)
	var comments []*Comment
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		var body *string
		err := rows.Scan(
			&comment.SenderName,
			&comment.Text,
			&comment.Received,
			&comment.Sender,
			&body,
			&comment.Pending,
			&comment.Subject,
			&comment.IsExtended,
			&comment.Id,
			&comment.ThreadHash,
		)
		if err != nil {
			return nil, err
		}
		if body != nil {
			comment.Body = *body
		}
		if comment.Subject == "" {
			comment.Subject = "(no subject)"
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (d *DAO) AllKeys() ([]*Key, error) {
	rows, err := d.Db.Query(`
        SELECT name, address, sigkey, deckey, enabled
        FROM keys2
        WHERE enabled=1
        ORDER BY name ASC
        LIMIT 10000;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := make([]*Key, 0)
	for rows.Next() {
		var name, address, sigkey, deckey string
		var enabled bool
		err := rows.Scan(&name, &address, &sigkey, &deckey, &enabled)
		if err != nil {
			return nil, err
		}
		if len(address) < 3 || address[:3] != "BM-" {
			address = "BM-" + address
		}
		keys = append(keys, &Key{
			Name:    name,
			Address: address,
			Sigkey:  sigkey,
			Deckey:  deckey,
			Enabled: enabled,
		})
	}
	return keys, nil
}

func (d *DAO) KeyByName(name string) *Key {
	key := Key{
		Name: name,
	}
	_ = d.Db.QueryRow(`
		select address, sigkey, deckey, enabled
		from keys2
		where name = ?;
	`, key.Name).Scan(&key.Address, &key.Sigkey, &key.Deckey, &key.Enabled)
	if key.Address != "" && (len(key.Address) < 3 || key.Address[:3] != "BM-") {
		key.Address = "BM-" + key.Address
	}
	return &key
}

func (d *DAO) KeyByName2(name string) (*Key, error) {
	key := Key{
		Name: name,
	}
	err := d.Db.QueryRow(`
		select address, sigkey, deckey, enabled
		from keys2
		where name = ?;
	`, key.Name).Scan(&key.Address, &key.Sigkey, &key.Deckey, &key.Enabled)
	if err != nil {
		return nil, err
	}
	if key.Address != "" && (len(key.Address) < 3 || key.Address[:3] != "BM-") {
		key.Address = "BM-" + key.Address
	}
	return &key, nil
}

func (d *DAO) ThreadsByChanName(name string) ([]*Thread, error) {
	rows, err := d.Db.Query(`
        select subject, count, last, hash
        from threads
        where name = ? and last > unix_timestamp() - 28*24*60*60
        order by last desc
        limit 100;
    `, "[chan] "+name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	threads := make([]*Thread, 0)
	for rows.Next() {
		var thread Thread
		err := rows.Scan(
			&thread.Subject,
			&thread.Count,
			&thread.Last,
			&thread.Hash,
		)
		if err != nil {
			return nil, err
		}
		if thread.Subject == "" {
			thread.Subject = "(no subject)"
		}
		threads = append(threads, &thread)
	}
	return threads, nil
}

func (d *DAO) TopChans() ([]*Chan, error) {
	chans := make([]*Chan, 0)
	rows, err := d.Db.Query(`
		select weekly, topchannels.name, last
		from topchannels join keys2
		on substr(topchannels.name, 8)=keys2.name
		order by spam asc, weekly desc, last desc
		limit 200;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		var name string
		var last int64
		err := rows.Scan(&count, &name, &last)
		if err != nil {
			return nil, err
		}
		chans = append(chans, &Chan{
			count,
			name[7:],
			last,
		})
	}
	return chans, nil
}
