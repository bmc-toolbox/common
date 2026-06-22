package common

import "testing"

func TestFormatVendorName(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		// exact-match aliases
		{"hpe alias", "HPE", VendorHPE},
		{"hp alias", "hp", VendorHPE},
		{"ami alias", "AMI", VendorAmericanMegatrends},

		// substring matches
		{"dell", "Dell Inc.", VendorDell},
		{"supermicro", "Supermicro", VendorSupermicro},
		{"lenovo", "Lenovo", VendorLenovo},
		{"lenovo lower", "lenovo", VendorLenovo},
		{"lenovo with text", "Lenovo Group Limited", VendorLenovo},
		{"lenovo with whitespace", "  Lenovo  ", VendorLenovo},
		{"marvell", "Marvell", VendorMarvell},

		// unknown vendor is returned unchanged
		{"unknown", "Acme Corp", "Acme Corp"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := FormatVendorName(tc.in); got != tc.want {
				t.Errorf("FormatVendorName(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatProductName(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"lenovo sr630 v2", "ThinkSystem SR630 V2", "sr630v2"},
		{"lenovo sr650", "ThinkSystem SR650", "sr650"},

		// unknown product name is returned unchanged
		{"unknown", "ThinkSystem SR999", "ThinkSystem SR999"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := FormatProductName(tc.in); got != tc.want {
				t.Errorf("FormatProductName(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
