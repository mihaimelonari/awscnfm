package release

import (
	"strconv"
	"strings"
	"testing"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/release/v1alpha1"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Release_mustFindMinor(t *testing.T) {
	testCases := []struct {
		name               string
		version            string
		release            v1alpha1.Release
		releases           []v1alpha1.Release
		expectedPrevious   v1alpha1.Release
		expectedLatest     v1alpha1.Release
		expectedUpgradable bool
	}{
		{
			name:    "case 0",
			version: "v12.0.0",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
			},
			expectedPrevious: v1alpha1.Release{},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"},
			},
			expectedUpgradable: false,
		},
		{
			name:    "case 1",
			version: "v12.0.0",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.1"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.2"}},
			},
			expectedPrevious: v1alpha1.Release{},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.0.2"},
			},
			expectedUpgradable: false,
		},
		{
			name:    "case 2",
			version: "v12.2.5",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.2.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"}},
			},
			expectedPrevious: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"},
			},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.2.5"},
			},
			expectedUpgradable: true,
		},
		{
			name:    "case 3",
			version: "v12.1.3",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"}},
			},
			expectedPrevious: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"},
			},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"},
			},
			expectedUpgradable: true,
		},
		{
			name:    "case 4",
			version: "v12.7.0",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"}},
			},
			expectedPrevious: v1alpha1.Release{},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"},
			},
			expectedUpgradable: false,
		},
		{
			name:    "case 5",
			version: "v12.1.1-dev",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.2-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1-dev"}},
			},
			expectedPrevious: v1alpha1.Release{},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"},
			},
			expectedUpgradable: false,
		},
		{
			name:    "case 6",
			version: "v12.4.4-dev",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.2-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1-dev"}},
			},
			expectedPrevious: v1alpha1.Release{},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"},
			},
			expectedUpgradable: false,
		},
		{
			name:    "case 7",
			version: "v12.2.3-dev",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.2.0-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.7-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.3-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.2-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.2.3-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"}},
			},
			expectedPrevious: v1alpha1.Release{},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"},
			},
			expectedUpgradable: false,
		},
		{
			name:    "case 8",
			version: "v100.1.6-xh3b4sd",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0-dev"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v100.1.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v100.0.3-xh3b4sd"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v100.1.3"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"}},
			},
			expectedPrevious: v1alpha1.Release{},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v100.0.3-xh3b4sd"},
			},
			expectedUpgradable: false,
		},
		{
			name:    "case 9",
			version: "v12.2.3-dev",
			releases: []v1alpha1.Release{
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.2.0"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.7"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.3"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.2"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.2.3"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.0.5"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "v12.1.1"}},
			},
			expectedPrevious: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.1.7"},
			},
			expectedLatest: v1alpha1.Release{
				ObjectMeta: metav1.ObjectMeta{Name: "v12.2.3"},
			},
			expectedUpgradable: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var r Resolver
			{
				c := MinorConfig{
					FromProject: tc.version,
					Releases:    tc.releases,
				}

				r, err = NewMinor(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			v := r.Version()

			expectedPrevious := strings.Replace(tc.expectedPrevious.GetName(), "v", "", 1)
			expectedLatest := strings.Replace(tc.expectedLatest.GetName(), "v", "", 1)

			if !cmp.Equal(v.Previous(), expectedPrevious) {
				t.Fatalf("\n\n%s\n", cmp.Diff(expectedPrevious, v.Previous()))
			}
			if !cmp.Equal(v.Latest(), expectedLatest) {
				t.Fatalf("\n\n%s\n", cmp.Diff(expectedLatest, v.Latest()))
			}
			if !cmp.Equal(r.Upgradable(), tc.expectedUpgradable) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.expectedUpgradable, r.Upgradable()))
			}
		})
	}
}