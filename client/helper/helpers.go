package helper

import (
	"client/config"
	"client/models"
)

func SyncPeers(newPeers []models.Peer) ([]models.Peer, []models.Peer) {
	added := []models.Peer{}
	removed := []models.Peer{}

	for _, peer := range newPeers {
		_, ok := config.Peers[peer.PublicKey]
		if ok == false {
			added = append(added, peer)
		}
	}
	for existingPeerPublicKey, existingPeer := range config.Peers {
		doesExist := false
		for _, newPeer := range newPeers {
			if newPeer.PublicKey == existingPeerPublicKey {
				doesExist = true
				break
			}
		}
		if doesExist == false {
			removed = append(removed, existingPeer)
		}

	}
	return added, removed
}
