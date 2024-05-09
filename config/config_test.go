package config

import "testing"

func Test_UnmarshalYamlConfig(t *testing.T) {
	c, err := New[ExtraData]("../test/config/data/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(c)

	cShort, err := New[ExtraData]("../test/config/data/config.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(cShort)
}

func Test_UnmarshalJsonConfig(t *testing.T) {
	c, err := New[ExtraData]("../test/config/data/config.json")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(c)
}
