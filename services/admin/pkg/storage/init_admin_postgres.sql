INSERT INTO public.categories(name) VALUES ('Молочные продукты');
INSERT INTO public.categories(name) VALUES ('Мясо и мясные изделия');
INSERT INTO public.categories(name) VALUES ('Рыба и морепродукты');
INSERT INTO public.categories(name) VALUES ('Овощи и фрукты');
INSERT INTO public.categories(name) VALUES ('Бакалея (крупы, макароны, мука)');
INSERT INTO public.categories(name) VALUES ('Хлеб и выпечка');
INSERT INTO public.categories(name) VALUES ('Кондитерские изделия');
INSERT INTO public.categories(name) VALUES ('Напитки (соки, воды, газированные напитки)');
INSERT INTO public.categories(name) VALUES ('Бытовая химия');
INSERT INTO public.categories(name) VALUES ('Товары для личной гигиены (средства для ухода за телом, косметика)');


-- Молочные продукты
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Молоко 2.5% жирности', 10, 50.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Сметана 15% жирности', 8, 75.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Творог 5% жирности', 5, 60.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Йогурт без добавок', 12, 45.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Кефир 1% жирности', 9, 55.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Сливки 10% жирности', 6, 80.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Сыр твёрдый «Голландский»', 4, 120.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Масло сливочное 82% жирности', 7, 90.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Простокваша', 11, 52.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (1, 'Ряженка', 13, 58.00);


-- Мясо и мясные изделия
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Говядина вырезка', 3, 350.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Свинина лопатка', 5, 280.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Курица бройлер', 8, 200.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Фарш мясной (говядина и свинина)', 6, 250.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Колбаса варёная', 10, 220.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Сардельки', 9, 180.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Ветчина', 4, 300.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Шницель свиной', 7, 270.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Филе индейки', 6, 260.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (2, 'Рулька свиная', 5, 290.00);


-- Рыба и морепродукты
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Лосось свежий', 4, 400.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Треска филе', 6, 320.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Сёмга копчёная', 3, 450.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Минтай замороженный', 8, 240.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Креветки варёные', 10, 380.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Кальмары замороженные', 7, 310.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Сельдь солёная', 9, 210.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Горбуша холодного копчения', 5, 360.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Камбала свежая', 6, 330.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (3, 'Икра красная', 2, 500.00);


-- Овощи и фрукты
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Помидоры свежие', 15, 80.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Огурцы свежие', 12, 70.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Морковь', 10, 65.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Лук репчатый', 8, 50.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Картофель', 20, 90.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Яблоки красные', 14, 95.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Бананы', 16, 100.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Апельсины', 13, 110.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Груши', 11, 105.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (4, 'Виноград', 9, 120.00);


-- Бакалея (крупы, макароны, мука)
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Рис белый круглозёрный', 25, 100.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Макароны из твёрдых сортов пшеницы', 20, 150.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Гречка ядрица', 18, 130.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Пшено', 16, 95.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Манная крупа', 14, 85.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Перловка', 12, 110.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Овсяные хлопья', 19, 140.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Кукурузная крупа', 17, 125.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Вермишель', 22, 160.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (5, 'Чечевица', 15, 170.00);


-- Хлеб и выпечка
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Хлеб пшеничный', 30, 40.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Батон нарезной', 28, 4);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Батон нарезной', 28, 45.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Хлеб ржаной', 25, 50.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Булочка с изюмом', 20, 60.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Сдоба', 15, 70.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Крекеры', 18, 55.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Хлебцы', 16, 65.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Пирожок с повидлом', 22, 75.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Тостовый хлеб', 30, 48.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (6, 'Хачапури', 14, 80.00);


INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Зефир', 12, 90.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Мармелад', 18, 75.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Вафли', 16, 80.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Карамельки', 25, 60.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Халва', 10, 120.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Шоколадные конфеты', 8, 150.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Торт «Наполеон»', 5, 200.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (7, 'Пряники', 14, 95.00);


-- Напитки (соки, воды, газированные напитки)
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Вода минеральная без газа', 50, 30.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Сок апельсиновый', 40, 50.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Газировка «Кола»', 45, 45.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Чай чёрный в пакетиках', 35, 60.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Кофе растворимый', 30, 70.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Компот из сухофруктов', 25, 55.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Кисель', 32, 50.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Лимонад', 42, 48.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Энергетический напиток', 28, 65.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (8, 'Вино столовое', 15, 150.00);


-- Бытовая химия
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Стиральный порошок', 20, 250.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Средство для мытья посуды', 25, 180.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Чистящее средство для ванной и туалета', 18, 220.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Кондиционер для белья', 15, 200.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Мыло хозяйственное', 30, 150.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Губки для мытья посуды', 50, 100.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Чистящий порошок', 22, 190.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Средство от ржавчины', 12, 230.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Спрей для стёкол', 28, 170.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (9, 'Универсальное чистящее средство', 24, 210.00);


-- Товары для личной гигиены (средства для ухода за телом, косметика)
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Зубная паста', 35, 120.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Зубная щётка', 40, 100.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Мыло туалетное', 50, 90.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Шампунь', 30, 130.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Гель для душа', 28, 140.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Дезодорант', 25, 160.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Туалетная бумага', 60, 80.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Носовые платки', 45, 95.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Средства для бритья', 22, 180.00);
INSERT INTO public.products(category_id, name, amount, cost) VALUES (10, 'Косметическое молочко для лица', 18, 200.00);
