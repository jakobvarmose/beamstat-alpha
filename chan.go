package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jakobvarmose/gobitmessage/crypt"
)

func getChanList() []ChannelInfo {
	list := make([]ChannelInfo, 0)
	rows, err := db.Query(`
		select weekly, topchannels.name, last
		from topchannels join keys2
		on substr(topchannels.name, 8)=keys2.name
		order by spam asc, weekly desc, last desc
		limit 200;
	`)
	if err != nil {
		panic(err)
	}
	i := 1
	for rows.Next() {
		var count int
		var name string
		var last int64
		err := rows.Scan(&count, &name, &last)
		if err != nil {
			panic(err)
		}
		list = append(list, ChannelInfo{
			i,
			count,
			name[7:],
			time.Unix(last, 0).Format("Jan 2 15:04"),
		})
		i++
	}
	rows.Close()
	return list
}

var chanAddresses = make(map[string]string)

func getChan(name string) (*Channel2, error) {
	address := ""
	exists := false
	_ = db.QueryRow(`
		select concat("BM-", address)
		from keys2
		where name = ?;
	`, name).Scan(&address)
	if address == "" {
		address = chanAddresses[name]
		if address == "" {
			ripe := crypt.DeterministicPrivateCombo(name).Ripe()
			address = crypt.Address{4, 1, ripe}.String()
			chanAddresses[name] = address
		}
	} else {
		exists = true
	}
	rows, err := db.Query(`
		select subject, count, last, hash
		from threads
		where name = ? and last > unix_timestamp() - 28*24*60*60
		order by last desc
		limit 100;
	`, "[chan] "+name)
	var threads []*Thread2
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var subject, hash string
		var count int
		var last int64
		rows.Scan(&subject, &count, &last, &hash)
		threads = append(threads, &Thread2{
			Subject: subject,
			Hash:    hash,
			Last:    time.Unix(last, 0).Format("Jan 2 15:04"),
			Count:   fmt.Sprintf("%d", count),
		})
	}
	rows.Close()
	channel := &Channel2{
		Name:    name,
		Address: address,
		Threads: threads,
		Exists:  exists,
	}
	return channel, nil
}
func handleChan(w http.ResponseWriter, req *http.Request, name string) {
	switch name {
	case "UPPERCASE_GENERAL":
		name = "GENERAL"
	case "aneki", "aneki/":
		name = "/aneki/"
	case "b", "b/":
		name = "/b/"
	case "pol", "pol/":
		name = "/pol/"
	case "tech", "tech/":
		name = "/tech/"
	}

	type ThreadInfo struct {
		Subject string
	}
	type Info struct {
		Channels []ChannelInfo
		Channel  *Channel2
		Thread   *ThreadInfo
	}
	channel, err := getChan(name)
	if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}
	info := Info{
		Channels: getChanList(),
		Channel:  channel,
		Thread:   &ThreadInfo{},
	}
	/*

		<?php
		$st->bindValue(1, $time, PDO::PARAM_INT);
		$st->execute();
		$i = 1;
		foreach ($st->fetchAll() as $row) {
		    $name = substr($row[1], 7);
			$last = date('M j H:i', $row[2]);

		    $i += 1;
		}
		?>
	*/

	if err := tmpl.ExecuteTemplate(w, "chan", info); err != nil {
		logrus.Error(err)
	}
}

/*

    $st = $db->prepare("select subject, count(*), max(received), max(coalesce(thread_hash, ''))
                        from channels
                        join objects
                        on channels.hash = objects.hash
                        where ? < received and name = ?
                        group by subject
                        order by max(received) desc
                        limit 100;");
    $st->bindValue(1, $time - 28*24*60*60, PDO::PARAM_INT);
    $st->bindValue(2, $channel_name, PDO::PARAM_STR);
    $st->execute();
    $i = 1;
    $found = false;
    foreach ($st->fetchAll() as $row) {
        echo "<tr>";
        $subject = htmlspecialchars($row[0]);
        $thread_hash = htmlspecialchars($row[3]);
        if ($thread_hash) {
            $name = htmlspecialchars(substr($channel_name, 7));
            if ($subject == '') {
                $subject = '&nbsp;';
            }
            echo "<td><a style=\"display:block\" href=\"/chan/$name/$thread_hash\">${subject}</a></td>";
        } else {
            echo "<td>${subject}</td>";
        }
        $last = date('M j H:i:s', $row[2]);
        echo "<td style=\"white-space: nowrap;\">${last}</td>";
        echo "<td class=\"text-right\">${row[1]}</td>";
        echo "</tr>";
        $i += 1;
        $found = true;
    }
if (!$found) {
    echo "<tr><td colspan=\"3\">&lt;Not Available&gt;</td></tr>";
}*/
