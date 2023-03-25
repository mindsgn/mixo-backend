import { Router } from 'express';
import * as stream from './../controller/stream.controller';
import multer from 'multer';

const storage = multer.diskStorage({
  filename: function (req, file, cb) {
    cb(null, file.originalname)
  },
  destination: function (req, file, cb) {
    cb(null, './uploads')
  },
})

const upload = multer({ storage })

const router = Router();
router.get("/", stream.getStreaming);
router.post("/post", stream.postStreaming);
router.get("/status", upload.any('file'), stream.getStatus);
router.get("/upload", stream.upload);

export default router;