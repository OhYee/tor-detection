# To add a new cell, type '# %%'
# To add a new markdown cell, type '# %% [markdown]'
# %% [markdown]
# # 识别 Tor 流量和非 Tor 流量
#
# 数据集地址 https://www.unb.ca/cic/datasets/tor.html
# %% [markdown]
# ## 数据读取
#
# 读入流量数据

# %%
import seaborn as sns
import matplotlib.pyplot as plt
import xgboost as xgb
from sklearn.model_selection import train_test_split
import sklearn
import pandas as pd
import numpy as np


data = pd.read_csv("./merged_5s.csv")

data.head(10)

# %% [markdown]
# ## 特征处理及数据转换
#
# 计算并只保留需要的特征

# %%
data = data[:-8]
data['Flow Bytes/s'] = data['Flow Bytes/s'].apply(float)
data['Flow Packets/s'] = data['Flow Packets/s'].apply(float)

temp = data.columns.to_list()
temp.remove("label")
temp.remove("Source IP")
temp.remove("Source Port")
temp.remove("Destination IP")
temp.remove("Destination Port")
temp.remove("Flow Duration")
temp.remove("Protocol")


x = data[temp]
y = data["label"].apply(lambda x: 0 if x == "nonTOR" else 1)

print(x)
print(y)

# %% [markdown]
# ## 数据集分割
#
# 将 pandas 的 DataFrame 转换为 sklearn 的 dataset

# %%


x_train, x_test, y_train, y_test = train_test_split(x, y, test_size=0.3)

print(x_train)
print(y_train)
print(x_test)
print(y_test)


# %%
x_train['Flow Bytes/s'] = x_train['Flow Bytes/s'].apply(lambda x: float(x))
x_train['Flow Packets/s'] = x_train['Flow Packets/s'].apply(lambda x: float(x))

x_train[x_train['Flow Bytes/s'].isnull()]

# %% [markdown]
# ## 训练模型并验证

# %%
x_train['Flow Bytes/s']


# %%
y_train


# %%
params = {
    'objective': 'binary:logistic',
    'eval_metric': 'auc',
    'learning_rate': 0.02,
    'max_depth': 8,
    'subsample': 0.9,
    'reg_lambda': 10,
    # 'tree_method': 'gpu_hist',
    'seed': 2021,
    #     'single_precision_histogram': False
    # 'deterministic_histogram': True
}

# params = {'tree_method': 'gpu_hist', 'max_depth': 8, 'alpha': 0,'num_leaves': 80,"seed":1024,
#               'gamma': 0, 'subsample': 1, 'scale_pos_weight': 1, 'learning_rate': 0.05,
#           'objective':'binary:logistic', 'eval_metric': ['error','auc']}
# model = xgb.Booster(model_file='xgb_75115.txt')
train_dmatrix = xgb.DMatrix(x_train, label=y_train)
valid_dmatrix = xgb.DMatrix(x_test, label=y_test)
# train_dmatrix = xgb.DMatrix(X_train, label=y_train)
# valid_dmatrix = xgb.DMatrix(X_test, label=y_test)
early_stopping = 30
verbose_eval = 100


model = xgb.train(
    params,
    train_dmatrix,
    evals=[(train_dmatrix, 'train'), (valid_dmatrix, 'valid')],
    verbose_eval=verbose_eval,
    num_boost_round=200,
    early_stopping_rounds=early_stopping,
)


# %%
# %%matplotlib
fig, ax = plt.subplots(figsize=(40, 40))
# plt.figure(figsize=(3000, 3000))
xgb.plot_importance(model,
                    height=0.5,
                    ax=ax)
plt.show()
