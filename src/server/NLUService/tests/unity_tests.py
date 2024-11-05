import pytest
from src.utils.preprocess import TextProcessor

################
## UNIT TESTS ##
################

# Test the TextProcessor class methods

# Atenttion: The vectorization step is the final method of the "TextProcessor" class.

# Fixture to create a TextProcessor instance
@pytest.fixture
def processor():
    return TextProcessor()

# Fixture to provide a sample text
@pytest.fixture
def sample_text():
    return "O que foi registrado sobre a emissão de boletos no APIMEC?"

# Test lowercase conversion
def test_lowercase_standardization(processor, sample_text):
    assert processor.lowercase_standardization(sample_text) == sample_text.lower()

# Test text tokenization
def test_tokenize_text(processor, sample_text):
    tokens = processor.tokenize_text(sample_text)
    assert tokens == ["O", "que", "foi", "registrado", "sobre", "a", "emissão", "de", "boletos", "no", "APIMEC", "?"]

# Test removal of stop words
def test_remove_stop_words(processor, sample_text):
    tokens = processor.tokenize_text(sample_text)
    filtered_tokens = processor.remove_stop_words(tokens)
    expected_filtered_tokens = ["registrado", "emissão", "boletos", "APIMEC"]
    assert all(token in filtered_tokens for token in expected_filtered_tokens)

# Test named entity recognition
def test_entity_identification(processor, sample_text):
    entities = processor.entity_identification(sample_text)
    expected_entities = [("APIMEC", "ORG")]
    assert any(ent in entities for ent in expected_entities)

# Test text normalization
def test_normalize_text(processor, sample_text):
    tokens = processor.tokenize_text(sample_text)
    normalized_tokens = processor.normalize_text(tokens)
    expected_normalized_tokens = ["O", "que", "foi", "registrado", "sobre", "a", "emissão", "de", "boletos", "no", "APIMEC", "?"]
    assert all(token in normalized_tokens for token in expected_normalized_tokens)

# Test lemmatization of tokens
def test_lemmatize(processor, sample_text):
    tokens = processor.tokenize_text(sample_text)
    lemmatized_tokens = processor.lemmatize(tokens)
    expected_lemmatized_tokens = ["o", "que", "ser", "registrar", "sobre", "o", "emissão", "de", "boleto", "em o", "APIMEC", "?"]
    assert all(lemma in lemmatized_tokens for lemma in expected_lemmatized_tokens)

# Test punctuation removal
def test_remove_punctuation(processor, sample_text):
    cleaned_text = processor.remove_punctuation(sample_text)
    assert cleaned_text == "O que foi registrado sobre a emissão de boletos no APIMEC"

