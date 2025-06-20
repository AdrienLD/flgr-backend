-- Créer la base allégée
ATTACH DATABASE 'lite_lyrics.sqlite3' AS lite;

-- Tables cibles
CREATE TABLE lite.lyrics (
  id INTEGER PRIMARY KEY,
  synced_lyrics TEXT NOT NULL
);

CREATE TABLE lite.tracks (
  id INTEGER PRIMARY KEY,
  name_lower TEXT NOT NULL,
  artist_name_lower TEXT NOT NULL,
  last_lyrics_id INTEGER NOT NULL
);

CREATE VIRTUAL TABLE lite.tracks_fts USING fts5(
  name_lower,
  artist_name_lower,
  content='tracks',
  content_rowid='id',
  tokenize = 'unicode61 remove_diacritics 2'
);

-- Copier uniquement les lyrics valides
INSERT INTO lite.lyrics(id, synced_lyrics)
SELECT id, synced_lyrics
FROM lyrics
WHERE synced_lyrics IS NOT NULL;

-- Copier uniquement les tracks dont les paroles sont valides
INSERT INTO lite.tracks(id, name_lower, artist_name_lower, last_lyrics_id)
SELECT t.id, t.name_lower, t.artist_name_lower, t.last_lyrics_id
FROM tracks t
JOIN lite.lyrics l ON l.id = t.last_lyrics_id;

-- Rebuild FTS
INSERT INTO lite.tracks_fts(rowid, name_lower, artist_name_lower)
SELECT id, name_lower, artist_name_lower FROM lite.tracks;

-- Déconnecter proprement
DETACH DATABASE lite;
