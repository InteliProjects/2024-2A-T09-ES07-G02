package main

import (
	"testing"
)

// TestCountOccurrences checks the occurrence count for a specific tag with multiple synonyms.
func TestCountOccurrences(t *testing.T) {
	// Create a root node with synonyms for "wegucci".
	root := NewTreeNode("mister", "clayton")
	_ = root.AddSubTag("wegucci", "we gucci", "go guy", "we", "gucci")

	// Count occurrences of "we gucci" in the node and its synonyms.
	count, _ := root.CountOccurrences("we gucci")

	// Verify that the count matches the expected value of 3 (including synonyms).
	if count != 3 {
		t.Errorf("expected count of 3, got %d", count)
	}
}

// TestSimpleSynonym checks the occurrence count for a single tag.
func TestSimpleSynonym(t *testing.T) {
	// Create a root node with a single tag "test".
	root := NewTreeNode("test")

	// Count occurrences of "this is a test".
	count, _ := root.CountOccurrences("this is a test")

	// Verify that the count matches the expected value of 1.
	if count != 1 {
		t.Errorf("expected count of 1, got %d", count)
	}
}

// TestMultipleSynonyms checks the occurrence count for multiple synonyms.
func TestMultipleSynonyms(t *testing.T) {
	// Create a root node with tags "test" and "exam".
	root := NewTreeNode("test", "exam")

	// Count occurrences of "this is a test exam".
	count, _ := root.CountOccurrences("this is a test exam")

	// Verify that the count matches the expected value of 2.
	if count != 2 {
		t.Errorf("expected count of 2, got %d", count)
	}
}

// TestSubtagSynonyms checks the occurrence count for a tag with subtags having synonyms.
func TestSubtagSynonyms(t *testing.T) {
	// Create a root node with a subtag "sub1" having synonyms "test" and "exam".
	root := NewTreeNode("root")
	_ = root.AddSubTag("sub1", "test", "exam")

	// Count occurrences of "this is a test exam" in the root node and its subtags.
	count, bestSubtag := root.CountOccurrences("this is a test exam")

	// Verify that the count matches the expected value of 2.
	if count != 2 {
		t.Errorf("expected count of 2, got %d", count)
	}
	// Verify that the best subtag is "sub1" (since both "test" and "exam" match).
	if bestSubtag != "sub1" {
		t.Errorf("expected best subtag 'sub1', got '%s'", bestSubtag)
	}
}

// TestMultipleSubtags checks the counting of occurrences with multiple subtags having different synonyms.
func TestMultipleSubtags(t *testing.T) {
	// Create a root node with subtags "sub1" and "sub2", each with different synonyms.
	root := NewTreeNode("root")
	sub1 := root.AddSubTag("sub1", "test")
	sub2 := root.AddSubTag("sub2", "exam", "quiz")

	// Count occurrences in subtag "sub1" for the content "test test".
	count1, _ := sub1.CountOccurrences("test test")
	if count1 != 2 {
		t.Errorf("expected count of 2, got %d", count1)
	}

	// Count occurrences in subtag "sub2" for the content "exam quiz".
	count2, _ := sub2.CountOccurrences("exam quiz")
	if count2 != 2 {
		t.Errorf("expected count of 2, got %d", count2)
	}

	// Count occurrences in the root node for the content "test test exam quiz".
	totalCount, bestSubtag := root.CountOccurrences("test test exam quiz")

	// Verify that the total count matches the expected value of 4.
	if totalCount != 4 {
		t.Errorf("expected total count of 4, got %d", totalCount)
	}
	// Verify that the best subtag is "sub1" (since "test" appears more frequently).
	if bestSubtag != "sub1" {
		t.Errorf("expected best subtag 'sub1', got '%s'", bestSubtag)
	}
}

// TestRepeatedSynonyms checks the occurrence count for repeated synonyms.
func TestRepeatedSynonyms(t *testing.T) {
	// Create a root node with synonyms "apple" and "fruit".
	root := NewTreeNode("apple", "fruit")

	// Count occurrences of "apple fruit apple fruit".
	count, _ := root.CountOccurrences("apple fruit apple fruit")

	// Verify that the count matches the expected value of 4 (each synonym appears twice).
	if count != 4 {
		t.Errorf("expected count of 4, got %d", count)
	}
}

// TestCompoundSynonyms checks the occurrence count for compound synonyms.
func TestCompoundSynonyms(t *testing.T) {
	// Create a root node with compound synonyms "fast food" and "junk food".
	root := NewTreeNode("fast food", "junk food")

	// Count occurrences of "I like fast food and junk food".
	count, _ := root.CountOccurrences("I like fast food and junk food")

	// Verify that the count matches the expected value of 2 (each compound synonym appears once).
	if count != 2 {
		t.Errorf("expected count of 2, got %d", count)
	}
}

// TestNoOccurrences checks the occurrence count when no synonyms are present in the content.
func TestNoOccurrences(t *testing.T) {
	// Create a root node with tags "test" and "exam".
	root := NewTreeNode("test", "exam")

	// Count occurrences in the content "this is just a trial".
	count, _ := root.CountOccurrences("this is just a trial")

	// Verify that the count matches the expected value of 0 (no matches found).
	if count != 0 {
		t.Errorf("expected count of 0, got %d", count)
	}
}

// TestSubtagWithoutSynonyms checks handling of a subtag with no synonyms.
func TestSubtagWithoutSynonyms(t *testing.T) {
	// Create a root node with a subtag "emptySubtag" but no synonyms.
	root := NewTreeNode("root")
	_ = root.AddSubTag("emptySubtag") // Subtag without synonyms

	// Count occurrences of "root" in the root node and its subtags.
	count, bestSubtag := root.CountOccurrences("root")

	// Verify that the count matches the expected value of 1 (only the root node matches).
	if count != 1 {
		t.Errorf("expected count of 1, got %d", count)
	}
	// Verify that the best subtag is an empty string (no synonyms in subtags).
	if bestSubtag != "" {
		t.Errorf("expected best subtag '', got '%s'", bestSubtag)
	}
}

// TestOverlappingSynonyms checks the occurrence count when synonyms overlap.
func TestOverlappingSynonyms(t *testing.T) {
	// Create a root node with overlapping synonyms "car", "auto", and "automobile".
	root := NewTreeNode("car", "auto", "automobile")

	// Count occurrences of "car auto automobile".
	count, _ := root.CountOccurrences("car auto automobile")

	// Verify that the count matches the expected value of 4 (each synonym appears once).
	if count != 4 {
		t.Errorf("expected count of 4, got %d", count)
	}
}

// TestSubtagWithMoreOccurrences checks counting of occurrences when a subtag has more matches.
func TestSubtagWithMoreOccurrences(t *testing.T) {
	// Create a root node with subtags "sub1" and "sub2". Subtag "sub2" has more occurrences.
	root := NewTreeNode("root")
	_ = root.AddSubTag("sub1", "test")
	_ = root.AddSubTag("sub2", "exam")

	// Count occurrences in the content "test exam exam".
	totalCount, bestSubtag := root.CountOccurrences("test exam exam")

	// Verify that the total count matches the expected value of 3 (subtag "sub2" has more occurrences).
	if totalCount != 3 {
		t.Errorf("expected total count of 3, got %d", totalCount)
	}
	// Verify that the best subtag is "sub2" (since "exam" appears more frequently).
	if bestSubtag != "sub2" {
		t.Errorf("expected best subtag 'sub2', got '%s'", bestSubtag)
	}
}
