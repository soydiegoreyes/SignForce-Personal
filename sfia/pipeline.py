# -*- coding: utf-8 -*-
import os
import tensorflow as tf
from tensorflow.keras import layers, models, metrics, backend as K
from tensorflow.keras.callbacks import EarlyStopping, ModelCheckpoint, ReduceLROnPlateau
import numpy as np
import pandas as pd
import random
import itertools

# --- 1. FUNCIONES DE UTILIDAD Y CARGA DE DATOS ---
# (Tu función original estaba bien, la mantenemos casi igual)
def separar_pares(img_dir, personas_permitidas, negativos_por_positivo=1):
    pares = []
    labels = []
    # Filtramos solo directorios válidos
    persona_imagen = {p: [os.path.join(img_dir, p, f) for f in os.listdir(os.path.join(img_dir, p)) if f.lower().endswith(('.png', '.jpg', '.jpeg'))] for p in personas_permitidas if os.path.isdir(os.path.join(img_dir, p))}
    # Limpiamos personas sin imágenes suficientes después del filtrado
    persona_imagen = {p: imgs for p, imgs in persona_imagen.items() if len(imgs) >= 1}
    personas = list(persona_imagen.keys())
    if len(personas) < 2:
        print("Error: No hay suficientes personas con imágenes para crear pares.")
        return np.array([]), np.array([])
    print(f"Generando pares para {len(personas)} personas...")
    for persona in personas:
        imagenes = persona_imagen[persona]
        # Si solo hay 1 imagen, no podemos hacer pares positivos, solo negativos
        if len(imagenes) < 2:
             # Generar algunos negativos para esta persona solitaria
             for _ in range(negativos_por_positivo * 2):
                otra_persona = random.choice([p for p in personas if p != persona])
                img_neg = random.choice(persona_imagen[otra_persona])
                pares.append([imagenes[0], img_neg])
                labels.append(0) # Diferentes
             continue
        # Pares Positivos (misma persona)
        # Usamos itertools.combinations para ser más eficientes si hay pocas fotos
        for img1, img2 in itertools.combinations(imagenes, 2):
            pares.append([img1, img2])
            labels.append(1) # Iguales
            # Pares Negativos asociados (personas diferentes)
            for _ in range(negativos_por_positivo):
                otra_persona = random.choice([p for p in personas if p != persona])
                img_neg = random.choice(persona_imagen[otra_persona])
                # Hacemos el par negativo con la primera imagen del par positivo actual
                pares.append([img1, img_neg])
                labels.append(0) # Diferentes
    print(f"Total de pares generados: {len(labels)}")
    # Mezclamos los pares y las etiquetas al unísono
    c = list(zip(pares, labels))
    random.shuffle(c)
    pares, labels = zip(*c)
    return np.array(pares), np.array(labels, dtype="float32")


def load_and_preprocess(img1_path, img2_path, input_shape):
    def process_img(path):
        img = tf.io.read_file(path)
        img = tf.image.decode_jpeg(img, channels=input_shape[2])
        # Usamos 'resize_with_pad' para no deformar la cara si no es cuadrada
        img = tf.image.resize_with_pad(img, input_shape[0], input_shape[1])
        img = img / 255.0  # Normalización a [0, 1]
        return img
    return (process_img(img1_path), process_img(img2_path))


# --- 2. IMPLEMENTACIÓN DE CONTRASTIVE LOSS Y DISTANCIA ---

def euclidean_distance(vects):
    """Calcula la distancia euclidiana entre dos vectores."""
    x, y = vects
    sum_square = K.sum(K.square(x - y), axis=1, keepdims=True)
    # Agregamos un pequeño epsilon para evitar la división por cero en el gradiente
    return K.sqrt(K.maximum(sum_square, K.epsilon()))


def contrastive_loss_with_margin(margin):
    """Crea la función de pérdida contrastiva con un margen específico."""
    def contrastive_loss(y_true, y_pred):
        # y_true: 1 si son la misma persona, 0 si son diferentes
        # y_pred: la distancia euclidiana calculada por la red
        # Si son la misma persona (y_true=1), queremos minimizar la distancia al cuadrado.
        # loss = y_true * y_pred^2
        square_pred = K.square(y_pred)
        # Si son diferentes (y_true=0), queremos que la distancia sea al menos 'margin'.
        # Si la distancia ya es mayor que el margen, la pérdida es 0.
        # loss = (1 - y_true) * max(margin - y_pred, 0)^2
        margin_square = K.square(K.maximum(margin - y_pred, 0))
        # Promedio de la pérdida en el lote
        return K.mean(y_true * square_pred + (1 - y_true) * margin_square)
    return contrastive_loss


# --- 3. ARQUITECTURA DE LA RED CON DATA AUGMENTATION ---
class L2Normalize(layers.Layer):
    def call(self, inputs):
        return tf.math.l2_normalize(inputs, axis=1)

def build_base_network(input_shape):
    """
    Extractor de características. Incluye aumento de datos al inicio.
    Devuelve un embedding normalizado L2.
    """
    inputs = layers.Input(shape=input_shape, name="input_image")
    # --- DATA AUGMENTATION (Solo activo durante el entrenamiento) ---
    # Para LFW, las caras están semi-alineadas, así que no seamos muy agresivos.
    x = layers.RandomFlip("horizontal", name="aug_flip")(inputs)
    x = layers.RandomRotation(factor=0.05, fill_mode="nearest", name="aug_rotate")(x) # +/- 5% rotación
    x = layers.RandomTranslation(height_factor=0.05, width_factor=0.05, fill_mode="nearest", name="aug_translate")(x)
    x = layers.RandomZoom(height_factor=0.05, name="aug_zoom")(x)
    # ---------------------------------------------------------------
    # Bloque Convolucional 1
    x = layers.Conv2D(64, (7, 7), padding="same", strides=(2,2), activation='relu', name='conv1')(x)
    x = layers.BatchNormalization(name='bn1')(x)
    x = layers.MaxPooling2D((2, 2), padding="same", name='pool1')(x)
    x = layers.Dropout(0.1)(x) # Dropout ligero para evitar sobreajuste temprano
    # Bloque Convolucional 2
    x = layers.Conv2D(128, (5, 5), padding="same", activation='relu', name='conv2')(x)
    x = layers.BatchNormalization(name='bn2')(x)
    x = layers.MaxPooling2D((2, 2), padding="same", name='pool2')(x)
    x = layers.Dropout(0.2)(x)
    # Bloque Convolucional 3
    x = layers.Conv2D(256, (3, 3), padding="same", activation='relu', name='conv3')(x)
    x = layers.BatchNormalization(name='bn3')(x)
    x = layers.MaxPooling2D((2, 2), padding="same", name='pool3')(x)
    x = layers.Dropout(0.3)(x)
    # Bloque Convolucional 4 (Más profundo para capturar detalles complejos)
    x = layers.Conv2D(512, (3, 3), padding="same", activation='relu', name='conv4')(x)
    x = layers.BatchNormalization(name='bn4')(x)
    x = layers.GlobalAveragePooling2D(name='gap')(x) # Mejor que Flatten para caras no centradas
    # Capa Densa final (Embedding)
    x = layers.Dense(256, activation=None, name='fc_embedding')(x)
    # Normalización L2: Crucial para que la distancia euclidiana y contrastive loss funcionen bien.
    # Obliga a que todos los vectores tengan longitud 1, viviendo en la superficie de una hiperesfera.
    #x = layers.Lambda(lambda t: tf.math.l2_normalize(t, axis=1), name="l2_norm")(x)
    x = L2Normalize(name="l2_norm")(x)
    model = models.Model(inputs, x, name="base_feature_extractor")
    return model


# --- 4. PIPELINE PRINCIPAL ---

# Configuración
IMG_DIR = "/home/faces_dataset/lfw-deepfunneled/lfw-deepfunneled/" # ¡Asegúrate que esta ruta es correcta!
INPUT_SHAPE = (105, 105, 3)
BATCH_SIZE = 64 # Aumentado para GPU
EPOCHS = 100 # ¡Más épocas!
MARGIN = 1.0 # Margen para la Contrastive Loss

# Preparación de datos
print("Leyendo directorio de imágenes...")
personas_todas = [p for p in os.listdir(IMG_DIR) if os.path.isdir(os.path.join(IMG_DIR, p))]
random.shuffle(personas_todas)

# Usar un split 80/20 para tener más datos de entrenamiento dado que LFW es difícil
split = int(0.75 * len(personas_todas))
train_ids = personas_todas[:split]
val_ids = personas_todas[split:]

print(f"Personas en entrenamiento: {len(train_ids)}, Personas en validación: {len(val_ids)}")

# Generar pares. Aumentamos negativos_por_positivo para balancear mejor.
train_pairs, train_labels = separar_pares(IMG_DIR, train_ids, negativos_por_positivo=3)
val_pairs, val_labels = separar_pares(IMG_DIR, val_ids, negativos_por_positivo=3)

if len(train_pairs) == 0 or len(val_pairs) == 0:
    print("Error crítico: No se pudieron generar suficientes pares para entrenar.")
    exit()

# Crear Datasets de TensorFlow
train_ds = tf.data.Dataset.from_tensor_slices(((train_pairs[:,0], train_pairs[:,1]), train_labels))
train_ds = train_ds.map(lambda x, y: (load_and_preprocess(x[0], x[1], INPUT_SHAPE), y), num_parallel_calls=tf.data.AUTOTUNE)
train_ds = train_ds.shuffle(buffer_size=2048).batch(BATCH_SIZE).prefetch(tf.data.AUTOTUNE)

val_ds = tf.data.Dataset.from_tensor_slices(((val_pairs[:,0], val_pairs[:,1]), val_labels))
val_ds = val_ds.map(lambda x, y: (load_and_preprocess(x[0], x[1], INPUT_SHAPE), y), num_parallel_calls=tf.data.AUTOTUNE)
val_ds = val_ds.batch(BATCH_SIZE).prefetch(tf.data.AUTOTUNE)

# --- 5. CONSTRUCCIÓN Y ENTRENAMIENTO DEL MODELO SIAMÉS ---

# 1. Instanciar la red base (con aumento de datos incluido)
base_cnn = build_base_network(INPUT_SHAPE)

# 2. Definir las entradas
input_a = layers.Input(shape=INPUT_SHAPE, name="input_a")
input_b = layers.Input(shape=INPUT_SHAPE, name="input_b")

# 3. Pasar ambas imágenes por la misma red para obtener embeddings
embedding_a = base_cnn(input_a)
embedding_b = base_cnn(input_b)

# 4. Calcular la distancia entre los embeddings
# El output del modelo NO es una clasificación (0 o 1), sino la DISTANCIA.
distance = layers.Lambda(euclidean_distance, output_shape=(1,), name="distance_layer")([embedding_a, embedding_b])

# 5. Modelo final
siamese_model = models.Model(inputs=[input_a, input_b], outputs=distance, name="siamese_network")

# Compilación
# Usamos la función de pérdida personalizada.
# Nota: Accuracy no es una buena métrica aquí porque el output es una distancia continua, no una clase.
# Usaremos la pérdida de validación como principal indicador.
siamese_model.compile(loss=contrastive_loss_with_margin(margin=MARGIN), optimizer=tf.keras.optimizers.Adam(learning_rate=0.0005))

siamese_model.summary()

# Callbacks
callbacks = [EarlyStopping(patience=15, restore_best_weights=True, monitor='val_loss', verbose=1),ModelCheckpoint('mejor_modelo_siames_lfw.keras', save_best_only=True, monitor='val_loss', verbose=1),ReduceLROnPlateau(factor=0.5, patience=5, min_lr=0.000001, monitor='val_loss', verbose=1)]

print("\nIniciando entrenamiento con Data Augmentation y Contrastive Loss...")
# Entrenar
history = siamese_model.fit(train_ds,validation_data=val_ds,epochs=EPOCHS,callbacks=callbacks)

print("Entrenamiento finalizado.")
siamese_model.save("modelo_siames.keras")
siamese_model.save("modelo_siames.h5", save_format="tf")

base_cnn = siamese_model.get_layer("base_feature_extractor")
base_cnn.save("face_embedding_model.keras")
base_cnn.save("face_embedding_model.h5", save_format="tf")

base_cnn.save_weights("pesos_siames.weights.h5")

df=pd.DataFrame(history.history)
df.to_csv("history_model.csv",index=False)

# USA ESTA VERSIÓN QUE TIENE LOS NOMBRES ORIGINALES
def build_extractor_para_js(input_shape):
    inputs = layers.Input(shape=input_shape, name="input_image")
    # Bloque 1 - Los nombres deben ser IDENTICOS a los del entrenamiento
    x = layers.Conv2D(64, (7, 7), padding="same", strides=(2,2), activation='relu', name='conv1')(inputs)
    x = layers.BatchNormalization(name='bn1')(x)
    x = layers.MaxPooling2D((2, 2), padding="same", name='pool1')(x)
    # Bloque 2
    x = layers.Conv2D(128, (5, 5), padding="same", activation='relu', name='conv2')(x)
    x = layers.BatchNormalization(name='bn2')(x)
    x = layers.MaxPooling2D((2, 2), padding="same", name='pool2')(x)
    # Bloque 3
    x = layers.Conv2D(256, (3, 3), padding="same", activation='relu', name='conv3')(x)
    x = layers.BatchNormalization(name='bn3')(x)
    x = layers.MaxPooling2D((2, 2), padding="same", name='pool3')(x)
    # Bloque 4
    x = layers.Conv2D(512, (3, 3), padding="same", activation='relu', name='conv4')(x)
    x = layers.BatchNormalization(name='bn4')(x)
    x = layers.GlobalAveragePooling2D(name='gap')(x)
    x = layers.Dense(256, activation=None, name='fc_embedding')(x)
    x = L2Normalize(name="l2_norm")(x)
    return models.Model(inputs, x, name="base_feature_extractor")



extractor = build_extractor_para_js((105, 105, 3))
extractor.save("face_embedding_tfjs.keras")
weights_path = 'pesos_siames.weights.h5'
try:
    extractor.load_weights(weights_path, skip_mismatch=False)
    print("✅ ¡Pesos cargados exitosamente por nombre!")
except Exception as e:
    print("⚠️ Fallo la carga directa. Intentando modo permisivo...")
    extractor.load_weights(weights_path, by_name=True, skip_mismatch=True)
    print("✅ Pesos cargados (se omitieron capas no coincidentes).")

extractor.save("face_embedding_tfjs.keras")
extractor.export("face_embedding_savedmodel") # este es el bueno

#tensorflowjs_converter --input_format=tf_saved_model face_embedding_savedmodel ./facevector # este es el bueno

#tensorflowjs_converter --input_format=keras face_embedding_tfjs.keras ./facevector
#tensorflowjs_converter --input_format=keras face_embedding_tfjs.keras ./facevector
#tensorflowjs_converter --input_format=tf_saved_model --output_format=tfjs_layers_model ./modelo_siames.keras ./facevector
#tensorflowjs_converter --input_format=keras siames_model.h5 ./facevector
#tensorflowjs_converter --input_format=tf_saved_model --output_format=tfjs_layers_model face_model_fixed   ./documentFlow/facevector