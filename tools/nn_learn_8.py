import tensorflow as tf

from tensorflow import keras
import numpy as np
import pandas as pd
from nn_input_with_phase_and_scale import SIZE

df = pd.read_pickle('nn_input_with_phase_and_scale.gzip')

print('Loaded')

K = 0.0062531601041799952
mse = keras.losses.MeanSquaredError()
def sigmoid_loss_with_static(y_true, y_pred):
	global mse
	middles = y_pred[:, 0]
	ends = y_pred[:, 1]
	results = y_true[:, 0]
	statics = y_true[:, 1]
	phases = y_true[:, 2]
	scales = y_true[:, 3]
	size = results.shape[0]
	results_column = tf.reshape(results, (size, 1))
	statics_column = tf.reshape(statics, (size, 1))
	phase_column = tf.reshape(phases, (size, 1))
	scale_column = tf.reshape(scales, (size, 1))
	middle_column = tf.reshape(middles, (size, 1))
	end_column = tf.reshape(ends, (size, 1))
	return mse(results_column, tf.math.sigmoid(K * ((tf.multiply(middle_column, (256 - phase_column)) + tf.multiply(tf.multiply(end_column, phase_column), scale_column)) / 256 + statics_column)))

# Yeah this line is needed
tf.config.run_functions_eagerly(True)

# tf.data.experimental.enable_debug_mode()

model = keras.Sequential([
	keras.layers.InputLayer(input_shape=(SIZE,), dtype=bool),
	keras.layers.Dense(8, activation='relu', dtype='float32'),
	keras.layers.Dense(2, dtype='float32')
])
keras.utils.plot_model(model)
model.compile(loss=sigmoid_loss_with_static, optimizer=keras.optimizers.Adam())
model.summary()

history = model.fit(x = np.array(df['position'].tolist()), y=df[['result', 'static', 'phase', 'scale']].values, epochs=26)

model.save('nn_result_8')
