from pycaret.regression import load_model as pycaret_load
import keras


def load_semantic_regression_model():
    return pycaret_load('models/pkl/semantic_regression_lgbm')


def load_resource_records_regression_model():
    return pycaret_load('models/pkl/resourse_records_regression_lgbm_22')


def load_keras_lstm_model():
    return keras.models.load_model("models/keras/dga_lstm.h5")
