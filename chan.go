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
	topChans, err := d.TopChans()
	if err != nil {
		//FIXME
		return nil
	}
	i := 1
	for _, topChan := range topChans {
		list = append(list, ChannelInfo{
			i,
			topChan.Count,
			topChan.Name,
			time.Unix(topChan.Last, 0).Format("Jan 2 15:04"),
		})
		i++
	}
	return list
}

var chanAddresses = make(map[string]string)

func getChan(name string) (*Channel2, error) {
	exists := false
	key := d.KeyByName(name)
	if key.Address == "" {
		key.Address = chanAddresses[name]
		if key.Address == "" {
			ripe := crypt.DeterministicPrivateCombo(name).Ripe()
			key.Address = crypt.Address{
				Version: 4,
				Stream:  1,
				Ripe:    ripe,
			}.String()
			chanAddresses[name] = key.Address
		}
	} else {
		exists = true
	}
	threads1, err := d.ThreadsByChanName(name)
	if err != nil {
		return nil, err
	}
	var threads []*Thread2
	for _, thread := range threads1 {
		threads = append(threads, &Thread2{
			Subject: thread.Subject,
			Hash:    thread.Hash,
			Last:    time.Unix(thread.Last, 0).Format("Jan 2 15:04"),
			Count:   fmt.Sprintf("%d", thread.Count),
		})
	}
	channel := &Channel2{
		Name:    name,
		Address: key.Address,
		Threads: threads,
		Exists:  exists,
		Enabled: key.Enabled,
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
