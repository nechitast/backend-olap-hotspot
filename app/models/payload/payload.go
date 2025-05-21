package payload

type Location struct {
	Pulau     string
	Provinsi  string
	Kota      string
	Kecamatan string
}

type Time struct {
	Tahun    string
	Semester string
	Kuartal  string
	Bulan    string
	Hari     string
}

type Other struct {
	Confidence string
	Satelite   string
}
