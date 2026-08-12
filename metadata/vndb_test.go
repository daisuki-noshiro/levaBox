package metadata

import "testing"

func TestSearchVNDB(t *testing.T) {
	results, err := SearchVNDB("Summer Pockets")
	if err != nil {
		t.Fatal(err)
	}

	if len(results) == 0 {
		t.Fatal("VNDB 没有返回任何搜索结果")
	}

	for _, result := range results {
		t.Logf(
			"ID=%s 标题=%s 原标题=%s 发售=%s",
			result.ID,
			result.Title,
			result.AltTitle,
			result.Released,
		)

		for _, developer := range result.Developers {
			t.Logf("开发商：%s", developer.Name)
		}

		if result.Image != nil {
			t.Logf("封面：%s", result.Image.URL)
		}
	}
}
