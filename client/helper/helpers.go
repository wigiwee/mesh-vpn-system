package helper

// func SyncPeers(newPeers []models.Peer) ([]models.Peer, []models.Peer) {
// 	added := []models.Peer{}
// 	removed := []models.Peer{}

// 	for _, peer := range newPeers {
// 		_, ok := config.PeerState[peer.PublicKey]
// 		if ok == false {
// 			added = append(added, peer)
// 		}
// 	}

// 	for existingPeerPublicKey, existingPeerState := range config.PeerState {
// 		doesExist := false
// 		for _, newPeer := range newPeers {
// 			if newPeer.PublicKey == existingPeerPublicKey {
// 				doesExist = true
// 				break
// 			}
// 		}
// 		if doesExist == false {
// 			removed = append(removed, existingPeerState.Peer)
// 		}

// 	}
// 	return added, removed
// }

func checkPing() bool {
	return true
}
