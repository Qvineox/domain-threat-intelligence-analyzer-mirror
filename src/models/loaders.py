from pycaret.regression import load_model as pycaret_load
from keras.saving import load_model as keras_load

def load_semantic_regression_model():
    return pycaret_load('models/pkl/semantic_regression_lgbm')


def load_resource_records_regression_model():
    return pycaret_load('models/pkl/resourse_records_regression_lgbm_22')


def load_keras_lstm_model():
    return keras_load("models/keras/dga_lstm.keras")