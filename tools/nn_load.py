import tensorflow as tf

from tensorflow import keras
import numpy as np
import pandas as pd
from nn_input_with_phase_and_scale import SIZE, vectorize_fen, idxs_to_vector


K = 0.0062
mse = keras.losses.MeanSquaredError()
def sigmoid_loss_with_static(y_true, y_pred):
	global mse
	results = y_true[:, 0]
	statics = y_true[:, 1]
	size = results.shape[0]
	results_column = tf.reshape(results, (size, 1))
	statics_column = tf.reshape(statics, (size, 1))
	return mse(results_column, tf.math.sigmoid(K * (y_pred + statics_column)))

model = keras.models.load_model('nn_result_8', custom_objects={'sigmoid_loss_with_static': sigmoid_loss_with_static})

# print(model.predict(tf.reshape(idxs_to_vector(vectorize_fen('2r5/3b1pk1/3ppbp1/1p5p/1P1pP3/3P3P/R1PN1PP1/3B2K1 w - - 0 26'), SIZE), (1, SIZE))))
# print(model.predict(tf.reshape(idxs_to_vector(vectorize_fen('2rqr1k1/p5p1/5p1p/2pp1P2/P3n1bN/2P1Q3/1P2P1BP/3R1RK1 b - - 0 24'), SIZE), (1, SIZE))))

np.set_printoptions(threshold=np.inf)

for i in range(2):
	print(model.layers[i].get_weights()[0].shape)
	print(repr(model.layers[i].get_weights()[0]))
	print(model.layers[i].get_weights()[1].shape)
	print(repr(model.layers[i].get_weights()[1]))



