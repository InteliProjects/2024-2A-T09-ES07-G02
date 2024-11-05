import spacy
import string
import numpy as np

# Class to process text data using SpaCy and Word2Vec
class TextProcessor:
    # Constructor method
    def __init__(self):
        # Load the SpaCy model for Portuguese language processing
        self.nlp = spacy.load('pt_core_news_sm')

        # Dictionary to map slang/abbreviations to their full forms
        self.slang_dict = {
            "vc": "você", "blz": "beleza", "pq": "porque", "tb": "também",
            "msg": "mensagem", "dps": "depois", "n": "não", "q": "que",
            "kd": "cadê", "vdd": "verdade", "aki": "aqui", "flw": "falou",
            "td": "tudo", "tks": "obrigado", "bjs": "beijos", "pls": "por favor",
            "obg": "obrigado", "vlw": "valeu", "cmg": "comigo", "qdo": "quando",
            "pvc": "pode ver", "pvt": "privado"
        }

        # String of punctuation characters to be removed from the text
        self.punctuation = string.punctuation

    # Method to convert the entire text to lowercase
    def lowercase_standardization(self, text):
        """Convert the entire text to lowercase."""
        return text.lower()

    # Method to tokenize the text into words using SpaCy
    def tokenize_text(self, text):
        """Tokenize the text into words using SpaCy."""
        doc = self.nlp(text)
        return [token.text for token in doc]

    # Method to remove stop words from the tokenized list
    def remove_stop_words(self, tokens):
        """Remove stop words from the tokenized list."""
        return [token for token in tokens if not self.nlp.vocab[token].is_stop]

    # Method to identify named entities in the text using SpaCy
    def entity_identification(self, text):
        """Identify named entities in the text using SpaCy."""
        doc = self.nlp(text)
        return [(ent.text, ent.label_) for ent in doc.ents]

    # Method to replace slang and abbreviations with their full forms
    def normalize_text(self, tokens):
        """Replace slang and abbreviations with their full forms."""
        return [self.slang_dict.get(word, word) for word in tokens]

    # Method to lemmatize the tokens using SpaCy
    def lemmatize(self, tokens):
        """Lemmatize the tokens using SpaCy to get the base form of words."""
        doc = self.nlp(" ".join(tokens))
        return [token.lemma_ for token in doc]

    # Method to remove punctuation from the text
    def remove_punctuation(self, text):
        """Remove punctuation from the text."""
        return text.translate(str.maketrans('', '', self.punctuation))


    # Method to execute the complete text processing pipeline
    def process(self, text):
        """Execute the complete text processing pipeline."""
        # Remove punctuation from the text
        text = self.remove_punctuation(text)

        # Convert the text to lowercase
        text = self.lowercase_standardization(text)

        # Identify named entities before stop word removal
        entities = self.entity_identification(text)

        # Tokenize the text into words
        tokens = self.tokenize_text(text)

        # Normalize tokens by replacing slang/abbreviations
        tokens = self.normalize_text(tokens)

        # Remove stop words from the tokenized list
        tokens = self.remove_stop_words(tokens)

        # Lemmatize the tokens to their base forms
        lemmatized = self.lemmatize(tokens)

        # Return the processed data as a dictionary
        return {
            "tokens": tokens,
            "entities": entities,
            "lemmatized": lemmatized,
        }
