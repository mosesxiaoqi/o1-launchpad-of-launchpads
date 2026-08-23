package branding

import "testing"

func TestValidateBranding(t *testing.T) {
	if err := Validate("ai-pad", "AI Pad", "demo", "https://example.com/logo.png", "#12AbEF"); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ slug, name, description, logo, color string }{
		{slug: "AI-Pad", name: "AI"},
		{slug: "admin", name: "AI"},
		{slug: "ai-pad", name: ""},
		{slug: "ai-pad", name: "AI", description: string(make([]byte, 1001))},
		{slug: "ai-pad", name: "AI", logo: "http://example.com/logo.png"},
		{slug: "ai-pad", name: "AI", color: "red"},
	} {
		if err := Validate(test.slug, test.name, test.description, test.logo, test.color); err == nil {
			t.Fatalf("expected invalid branding: %#v", test)
		}
	}
}
