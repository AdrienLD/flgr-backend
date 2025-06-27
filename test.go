package main

type PlaylistItem struct {
	Titre   string `json:"titre"`
	Artiste string `json:"artiste"`
}

var Playlist60 = []PlaylistItem{
	{Titre: "Petit Rainbow", Artiste: "Sylvie Vartan"},
	{Titre: "Les p'tits papiers", Artiste: "Regine"},
	{Titre: "Debout les gars", Artiste: "Hugues Aufray"},
	{Titre: "San Francisco", Artiste: "Maxime Le Forestier"},
	{Titre: "Kiss Me", Artiste: "C. Jérôme"},
	{Titre: "Chanson inédite", Artiste: "Michèle Torr"},
	{Titre: "Madame", Artiste: "Claude Barzotti"},
	{Titre: "Elsa", Artiste: "Didier Barbelivien"},
	{Titre: "Tous les garçons et les filles - Stereo Mix", Artiste: "Françoise Hardy"},
	{Titre: "Mon amie la rose", Artiste: "Françoise Hardy"},
}

var Playlist70 = []PlaylistItem{
	{Titre: "YMCA - Original Version 1978", Artiste: "Village People"},
	{Titre: "Tu verras", Artiste: "Claude Nougaro"},
	{Titre: "Pas de boogie woogie", Artiste: "Eddy Mitchell"},
	{Titre: "Sex Machine", Artiste: "James Brown"},
	{Titre: "Bohemian Rhapsody - Remastered 2011", Artiste: "Queen"},
	{Titre: "Comme d'habitude", Artiste: "Claude François"},
	{Titre: "Mon frère", Artiste: "Maxime Le Forestier"},
	{Titre: "Le téléphone pleure", Artiste: "Claude François, Frédérique"},
	{Titre: "J'ai dix ans", Artiste: "Laurent Voulzy, Alain Souchon"},
	{Titre: "Waterloo", Artiste: "ABBA"},
}

var Playlist80 = []PlaylistItem{
	{Titre: "Sarà perché ti amo", Artiste: "Ricchi E Poveri"},
	{Titre: "Don't Go", Artiste: "Yazoo"},
	{Titre: "Trois nuits par semaine", Artiste: "Indochine"},
	{Titre: "Nuit de folie - Version originale 1988", Artiste: "Début De Soirée"},
	{Titre: "Africa", Artiste: "Rose Laurens"},
	{Titre: "I'm So Excited", Artiste: "The Pointer Sisters"},
	{Titre: "Ella, elle l'a - Remasterisé en 2004", Artiste: "France Gall"},
	{Titre: "J'ai encore rêvé d'elle", Artiste: "Il Etait Une Fois"},
	{Titre: "En rouge et noir", Artiste: "Jeanne Mas"},
	{Titre: "I Wanna Dance with Somebody (Who Loves Me)", Artiste: "Whitney Houston"},
}

var Playlist90 = []PlaylistItem{
	{Titre: "Mon papa à moi est un gangster", Artiste: "Stomy Bugsy"},
	{Titre: "I Don't Want to Miss a Thing - From \"Armageddon\" Soundtrack", Artiste: "Aerosmith"},
	{Titre: "Alane - Radio Version", Artiste: "Wes"},
	{Titre: "Doo Wop (That Thing)", Artiste: "Ms. Lauryn Hill"},
	{Titre: "No Diggity", Artiste: "Blackstreet, Dr. Dre, Queen Pen"},
	{Titre: "Don't Let the Sun Go Down on Me", Artiste: "George Michael, Elton John"},
	{Titre: "Enjoy the Silence", Artiste: "Depeche Mode"},
	{Titre: "Genie In a Bottle", Artiste: "Christina Aguilera"},
	{Titre: "Mr. Loverman (feat. Chevelle Franklin)", Artiste: "Shabba Ranks, Chevelle Franklin, David Morales, Hugo Dwyer"},
	{Titre: "No Ordinary Love", Artiste: "Sade"},
}

var Playlist00 = []PlaylistItem{
	{Titre: "Applause", Artiste: "Lady Gaga"},
	{Titre: "I Gotta Feeling - Edit", Artiste: "Black Eyed Peas"},
	{Titre: "Down To The River To Pray", Artiste: "Alison Krauss"},
	{Titre: "Club Can't Handle Me (feat. David Guetta)", Artiste: "Flo Rida, David Guetta"},
	{Titre: "Déconnectés", Artiste: "DJ Hamida, Lartiste, Kayna Samet, Rim'K"},
	{Titre: "Double je - Remix", Artiste: "Christophe Willem"},
	{Titre: "Can't Hold Us (feat. Ray Dalton)", Artiste: "Macklemore, Ryan Lewis, Macklemore & Ryan Lewis, Ray Dalton"},
	{Titre: "Hung Up", Artiste: "Madonna"},
	{Titre: "Evacuate The Dancefloor", Artiste: "Cascada"},
	{Titre: "Right Now (Na Na Na)", Artiste: "Akon"},
}

var Playlist10 = []PlaylistItem{
	{Titre: "Promiscuous", Artiste: "Nelly Furtado, Timbaland"},
	{Titre: "A Sky Full of Stars", Artiste: "Coldplay"},
	{Titre: "She Wolf (Falling to Pieces) [feat. Sia]", Artiste: "David Guetta, Sia"},
	{Titre: "Hey Baby (Drop It to the Floor) (feat. T-Pain) - Radio Edit", Artiste: "Pitbull, T-Pain"},
	{Titre: "Lonely", Artiste: "Akon"},
	{Titre: "Waka Waka (This Time for Africa) [The Official 2010 FIFA World Cup (TM) Song] (feat. Freshlyground)", Artiste: "Shakira, Freshlyground"},
	{Titre: "E.T.", Artiste: "Katy Perry"},
	{Titre: "Summer Jam 2003 - DJ F.R.A.N.K.'s Summermix Short", Artiste: "The Underdog Project"},
	{Titre: "Dernière danse", Artiste: "Kyo"},
	{Titre: "Poker Face", Artiste: "Lady Gaga"},
}

var AllPlaylists = map[string][]PlaylistItem{
	"60s": Playlist60,
	"70s": Playlist70,
	"80s": Playlist80,
	"90s": Playlist90,
	"00s": Playlist00,
	"10s": Playlist10,
}
