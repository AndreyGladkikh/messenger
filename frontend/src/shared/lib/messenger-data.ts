export type MessageStatus = 'sent' | 'delivered' | 'read'

export type Message = {
  id: string
  authorId: string
  text: string
  time: string
  outgoing: boolean
  status?: MessageStatus
  attachment?: {
    kind: 'file' | 'voice'
    name: string
    meta: string
  }
}

export type Conversation = {
  id: string
  name: string
  initials: string
  kind: 'direct' | 'group'
  online: boolean
  lastSeen: string
  members?: number
  preview: string
  previewOutgoing: boolean
  time: string
  unread: number
  muted: boolean
  pinned: boolean
  typing?: boolean
  messages: Message[]
}

export const currentUser = {
  name: 'Аня Соколова',
  initials: 'АС',
  handle: '@anya',
}

export const conversations: Conversation[] = [
  {
    id: 'design-team',
    name: 'Продуктовый дизайн',
    initials: 'ПД',
    kind: 'group',
    online: true,
    lastSeen: 'активны сейчас',
    members: 8,
    preview: 'Кирилл: Прикрепил финальную сетку, гляньте отступы',
    previewOutgoing: false,
    time: '14:32',
    unread: 3,
    muted: false,
    pinned: true,
    typing: true,
    messages: [
      {
        id: 'm1',
        authorId: 'Кирилл Ремизов',
        text: 'Собрал новую сетку для карточек диалогов — 12 колонок, но с уменьшенным гэпом на планшете.',
        time: '14:02',
        outgoing: false,
      },
      {
        id: 'm2',
        authorId: 'Кирилл Ремизов',
        text: 'Прикрепил финальную сетку, гляньте отступы',
        time: '14:03',
        outgoing: false,
        attachment: {
          kind: 'file',
          name: 'messenger-grid-v4.fig',
          meta: 'Figma · 4,2 МБ',
        },
      },
      {
        id: 'm3',
        authorId: currentUser.name,
        text: 'Отступы норм, но в списке чатов аватар лучше уменьшить до 40 — строка станет спокойнее.',
        time: '14:11',
        outgoing: true,
        status: 'read',
      },
      {
        id: 'm4',
        authorId: 'Лена Тихая',
        text: 'Поддерживаю. И давайте оставим один акцентный цвет, а статусы делать типографикой.',
        time: '14:18',
        outgoing: false,
      },
      {
        id: 'm5',
        authorId: currentUser.name,
        text: 'Да, ровно так и сделаем. Сегодня вечером закину прототип в ветку.',
        time: '14:26',
        outgoing: true,
        status: 'delivered',
      },
    ],
  },
  {
    id: 'marina',
    name: 'Марина Ковалёва',
    initials: 'МК',
    kind: 'direct',
    online: true,
    lastSeen: 'в сети',
    preview: 'Отправила правки по тексту на главной',
    previewOutgoing: false,
    time: '13:58',
    unread: 1,
    muted: false,
    pinned: true,
    messages: [
      {
        id: 'm1',
        authorId: 'Марина Ковалёва',
        text: 'Привет! Прошлась по тексту на главной, местами слишком длинные заголовки.',
        time: '13:40',
        outgoing: false,
      },
      {
        id: 'm2',
        authorId: currentUser.name,
        text: 'Спасибо! Скинь список, я поправлю вместе с версткой.',
        time: '13:52',
        outgoing: true,
        status: 'read',
      },
      {
        id: 'm3',
        authorId: 'Марина Ковалёва',
        text: 'Отправила правки по тексту на главной',
        time: '13:58',
        outgoing: false,
        attachment: {
          kind: 'file',
          name: 'copy-review.docx',
          meta: 'Документ · 128 КБ',
        },
      },
    ],
  },
  {
    id: 'ilya',
    name: 'Илья Дорохов',
    initials: 'ИД',
    kind: 'direct',
    online: false,
    lastSeen: 'был в сети 20 минут назад',
    preview: 'Вы: Созвон переносим на четверг',
    previewOutgoing: true,
    time: '12:14',
    unread: 0,
    muted: false,
    pinned: false,
    messages: [
      {
        id: 'm1',
        authorId: 'Илья Дорохов',
        text: 'По бэкенду успеваем, но нужна финальная схема сообщений.',
        time: '11:47',
        outgoing: false,
      },
      {
        id: 'm2',
        authorId: currentUser.name,
        text: 'Созвон переносим на четверг',
        time: '12:14',
        outgoing: true,
        status: 'read',
      },
    ],
  },
  {
    id: 'release',
    name: 'Релиз 2.4',
    initials: 'Р4',
    kind: 'group',
    online: false,
    lastSeen: 'последняя активность в 11:30',
    members: 14,
    preview: 'Соня: Сборка зелёная, можно катить на стейдж',
    previewOutgoing: false,
    time: '11:30',
    unread: 0,
    muted: true,
    pinned: false,
    messages: [
      {
        id: 'm1',
        authorId: 'Соня Ветрова',
        text: 'Сборка зелёная, можно катить на стейдж',
        time: '11:30',
        outgoing: false,
      },
    ],
  },
  {
    id: 'petr',
    name: 'Пётр Ким',
    initials: 'ПК',
    kind: 'direct',
    online: true,
    lastSeen: 'в сети',
    preview: 'Голосовое сообщение · 0:42',
    previewOutgoing: false,
    time: 'Вчера',
    unread: 0,
    muted: false,
    pinned: false,
    messages: [
      {
        id: 'm1',
        authorId: 'Пётр Ким',
        text: '',
        time: '19:05',
        outgoing: false,
        attachment: {
          kind: 'voice',
          name: 'Голосовое сообщение',
          meta: '0:42',
        },
      },
    ],
  },
  {
    id: 'family',
    name: 'Семья',
    initials: 'СМ',
    kind: 'group',
    online: false,
    lastSeen: 'последняя активность вчера',
    members: 5,
    preview: 'Мама: Не забудь про воскресенье',
    previewOutgoing: false,
    time: 'Вчера',
    unread: 0,
    muted: false,
    pinned: false,
    messages: [
      {
        id: 'm1',
        authorId: 'Мама',
        text: 'Не забудь про воскресенье',
        time: '18:20',
        outgoing: false,
      },
    ],
  },
  {
    id: 'saved',
    name: 'Избранное',
    initials: 'ИЗ',
    kind: 'direct',
    online: false,
    lastSeen: 'только вы',
    preview: 'Вы: Ссылка на статью про тайпографику',
    previewOutgoing: true,
    time: 'Пн',
    unread: 0,
    muted: false,
    pinned: false,
    messages: [
      {
        id: 'm1',
        authorId: currentUser.name,
        text: 'Ссылка на статью про тайпографику',
        time: '09:12',
        outgoing: true,
        status: 'sent',
      },
    ],
  },
]
